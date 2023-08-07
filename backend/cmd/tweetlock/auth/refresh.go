package main

import (
	"fmt"
	"serverless/internal/auth"
	"serverless/internal/requests"
	"serverless/pkg/tweetlock/datastructs"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func refreshHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reject := func() events.APIGatewayProxyResponse {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers: map[string]string{
				"Set-Cookie": requests.RefreshCookie("", time.Unix(0, 0)),
			},
		}
	}

	refreshToken := requests.GetCookieValue(r.Headers["Cookie"], "x-refresh")
	if refreshToken == "" {
		return reject(), nil
	}

	result, err := dynamo.Query(&dynamodb.QueryInput{
		TableName:              aws.String("tweetlock-sessions"),
		Limit:                  aws.Int64(1),
		KeyConditionExpression: aws.String("session_id = :session_id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":session_id": {S: aws.String(refreshToken)},
			":expires_at": {N: aws.String(fmt.Sprintf("%d", time.Now().Unix()))},
		},
		FilterExpression: aws.String("expires_at > :expires_at"),
	})
	if err != nil {
		return requests.BuildResponse(500, nil), err
	}
	if result == nil || *result.Count == 0 {
		return reject(), nil
	}

	session := datastructs.Session{}
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &session); err != nil {
		return requests.BuildResponse(500, nil), err
	}

	result, err = dynamo.Query(&dynamodb.QueryInput{
		TableName:              aws.String("tweetlock-users"),
		Limit:                  aws.Int64(1),
		KeyConditionExpression: aws.String("id = :id"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":id": {S: aws.String(session.UserID)},
		},
	})
	if err != nil {
		return requests.BuildResponse(500, nil), err
	}
	if result == nil || *result.Count == 0 {
		return reject(), nil
	}
	user := datastructs.User{}
	if err := dynamodbattribute.UnmarshalMap(result.Items[0], &user); err != nil {
		return requests.BuildResponse(500, nil), err
	}

	signed, expires, err := auth.SignAccessToken(user)
	if err != nil {
		return requests.BuildResponse(500, nil), err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Set-Cookie": requests.AccessCookie(signed, expires),
		},
	}, nil
}
