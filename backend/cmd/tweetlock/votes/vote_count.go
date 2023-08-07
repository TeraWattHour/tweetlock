package main

import (
	"fmt"
	"regexp"
	"serverless/internal/auth"
	"serverless/internal/requests"
	"serverless/pkg/tweetlock/datastructs"
	"serverless/pkg/tweetlock/dtos"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"golang.org/x/exp/slices"
)

var VALID_TARGETS_RE = regexp.MustCompile(`^([0-9]{1,19},){0,}([0-9]{1,19})$`)

func voteCountHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID, isAuth := auth.Guard(r)
	if !isAuth {
		return requests.BuildResponse(401, nil), nil
	}

	targets := r.QueryStringParameters["targets"]

	if !VALID_TARGETS_RE.MatchString(targets) {
		return requests.BuildResponse(422, map[string]interface{}{
			"message": "Invalid request. Provide a list of Twitter IDs separated by commas as a query parameter `targets`",
		}), nil
	}

	targetIds := strings.Split(targets, ",")
	uniqueIds := []string{}

	for _, value := range targetIds {
		if !slices.Contains(uniqueIds, value) {
			uniqueIds = append(uniqueIds, value)
		}
	}

	tempRows := ""
	for k, v := range uniqueIds {
		if k == 0 {
			tempRows += fmt.Sprintf(`select '%s' as twitter_id`, v)
		} else {
			tempRows += fmt.Sprintf(` select '%s'`, v)
		}
		if k != len(uniqueIds)-1 {
			tempRows += " union all"
		}
	}

	te := []datastructs.VoteCount{}

	query := fmt.Sprintf(`
		select t.*, count(v.target_id) as votes, (
			select count(*) = 1 from votes vo where vo.target_id = twitter_id and vo.user_id = ?
		) as has_voted from ( %s ) as t
		left join votes v on v.target_id = t.twitter_id
		group by t.twitter_id;
	`, tempRows)

	err := DB.Select(&te, query, userID)
	if err != nil {
		return requests.BuildResponse(500, nil), nil
	}

	result := dtos.VoteCountMap{}

	for _, entry := range te {
		result[entry.TwitterID] = dtos.VoteCount{
			Votes:    entry.Votes,
			HasVoted: entry.HasVoted,
		}
	}

	return requests.BuildResponse(200, result), nil
}
