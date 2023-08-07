package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-lambda-go/events"
)

func twitterRedirectHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", os.Getenv("TWITTER_CLIENT_ID"))
	params.Add("redirect_uri", os.Getenv("TWITTER_REDIRECT"))
	params.Add("state", "state")
	params.Add("scope", "offline.access tweet.read users.read")
	params.Add("code_challenge", "challenge")
	params.Add("code_challenge_method", "plain")

	return events.APIGatewayProxyResponse{
		StatusCode: 307,
		Headers: map[string]string{
			"Location": fmt.Sprintf("https://twitter.com/i/oauth2/authorize?%s", params.Encode()),
		},
	}, nil
}
