package dtos

type VoteCount struct {
	HasVoted bool `json:"hasVoted"`
	Votes    uint `json:"votes"`
}

type VoteCountMap map[string]VoteCount
