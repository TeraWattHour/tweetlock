package main

import (
	"errors"
	"serverless/internal/auth"
	"serverless/internal/requests"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-sql-driver/mysql"
)

func voteHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID, isAuth := auth.Guard(r)
	if !isAuth {
		return requests.BuildResponse(401, nil), nil
	}

	targetID := r.PathParameters["target"]

	if r.HTTPMethod == "POST" {
		code := 201

		_, err := DB.Exec(`insert into votes(target_id, user_id) values (?, ?)`, targetID, userID)
		if err != nil {
			var mysqlError *mysql.MySQLError
			if !errors.As(err, &mysqlError) || mysqlError.Number != 1062 {
				return requests.BuildResponse(500, nil), err
			}
			code = 200
		}

		return requests.BuildResponse(code, nil), nil
	}
	if r.HTTPMethod == "DELETE" {
		_, err := DB.Exec(`delete from votes where target_id = ? and user_id = ?`, targetID, userID)
		if err != nil {
			return requests.BuildResponse(500, nil), err
		}

		return requests.BuildResponse(200, nil), nil
	}

	return requests.BuildResponse(405, nil), nil
}
