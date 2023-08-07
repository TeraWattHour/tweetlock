package requests

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func BuildResponse(statusCode int, body interface{}) events.APIGatewayProxyResponse {
	ev := events.APIGatewayProxyResponse{StatusCode: statusCode}
	if body == nil {
		return ev
	}
	by, err := json.Marshal(body)

	if err != nil {
		ev.StatusCode = 500
		return ev
	}

	ev.Body = string(by)
	ev.Headers = map[string]string{
		"Content-Type": "application/json",
	}

	return ev
}
