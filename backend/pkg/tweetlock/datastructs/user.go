package datastructs

type User struct {
	ID        string `json:"id"`
	TwitterID string `dynamodbav:"twitter_id" json:"twitterId"`

	Name   string `json:"name"`
	Handle string `json:"twitterName"`
}
