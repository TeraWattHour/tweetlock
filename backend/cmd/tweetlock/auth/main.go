package main

import (
	"os"
	"serverless/internal/requests"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB = sqlx.MustConnect("mysql", os.Getenv("DB_DSN"))
var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var dynamo = dynamodb.New(sess)

func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if strings.HasPrefix(r.Path, "/twitter-redirect") && r.HTTPMethod == "GET" {
		return twitterRedirectHandler(r)
	}

	if strings.HasPrefix(r.Path, "/twitter-callback") && r.HTTPMethod == "GET" {
		return twitterCallbackHandler(r)
	}

	if strings.HasPrefix(r.Path, "/refresh") && r.HTTPMethod == "POST" {
		return refreshHandler(r)
	}

	return requests.BuildResponse(404, nil), nil
}

func main() {
	lambda.Start(handler)
}
