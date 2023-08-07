package main

import (
	"fmt"
	"serverless/internal/auth"
	"serverless/internal/requests"
	"serverless/internal/twitter"
	"serverless/pkg/tweetlock/datastructs"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

func twitterCallbackHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	_, isAuth := auth.Guard(r)
	if isAuth {
		return requests.BuildResponse(403, nil), nil
	}

	code := r.QueryStringParameters["code"]
	tokens, err := twitter.GetTokensFromCode(code)
	if err != nil {
		return requests.BuildResponse(400, nil), nil
	}

	twitterUser, err := twitter.GetUserData(tokens.AccessToken)
	if err != nil {
		return requests.BuildResponse(400, nil), nil
	}

	dbUser := datastructs.User{}
	result, err := dynamo.Scan(&dynamodb.ScanInput{
		TableName:        aws.String("tweetlock-users"),
		FilterExpression: aws.String("twitter_id = :twitter_id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":twitter_id": {S: aws.String(twitterUser.ID)},
		},
	})
	if err != nil {
		return requests.BuildResponse(500, nil), err
	}
	if result == nil || *result.Count == 0 {
		dbUser.ID = uuid.NewString()
		_, err = dynamo.PutItem(&dynamodb.PutItemInput{
			TableName: aws.String("tweetlock-users"),
			Item: map[string]*dynamodb.AttributeValue{
				"id":         {S: aws.String(dbUser.ID)},
				"twitter_id": {S: aws.String(twitterUser.ID)},
				"name":       {S: aws.String(twitterUser.Name)},
				"handle":     {S: aws.String(twitterUser.Handle)},
			},
		})
		if err != nil {
			return requests.BuildResponse(500, nil), err
		}
		dbUser.TwitterID = twitterUser.ID
		dbUser.Name = twitterUser.Name
		dbUser.Handle = twitterUser.Handle
	} else {
		err = dynamodbattribute.UnmarshalMap(result.Items[0], &dbUser)
		if err != nil {
			return requests.BuildResponse(500, nil), err
		}
	}

	sessionId := uuid.NewString()
	refreshExpires := time.Now().Add(time.Hour * 24 * 365)
	_, err = dynamo.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("tweetlock-sessions"),
		Item: map[string]*dynamodb.AttributeValue{
			"session_id": {S: aws.String(sessionId)},
			"user_id":    {S: aws.String(dbUser.ID)},
			"expires_at": {N: aws.String(fmt.Sprintf("%d", refreshExpires.Unix()))},
		},
	})
	if err != nil {
		return requests.BuildResponse(500, nil), err
	}

	accessToken, accessExpires, err := auth.SignAccessToken(dbUser)
	if err != nil {
		return requests.BuildResponse(500, nil), err
	}

	accessCookie := requests.AccessCookie(accessToken, accessExpires)
	refreshCookie := requests.RefreshCookie(sessionId, refreshExpires)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "<h1>Successfully signed in! You may now close this page.</h1>",
		Headers: map[string]string{
			"Content-Type": "text/html",
		},
		MultiValueHeaders: map[string][]string{
			"Set-Cookie": {accessCookie, refreshCookie},
		},
	}, nil
}
