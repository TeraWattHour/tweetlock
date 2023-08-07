package main

import (
	"os"
	"serverless/internal/requests"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DB = sqlx.MustConnect("mysql", os.Getenv("DB_DSN"))

func handler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if strings.HasPrefix(r.Path, "/vote-count") && r.HTTPMethod == "GET" {
		return voteCountHandler(r)
	}

	if strings.HasPrefix(r.Path, "/vote") {
		return voteHandler(r)
	}

	return requests.BuildResponse(404, nil), nil
}

func main() {
	lambda.Start(handler)
}
