package dtos

type User struct {
	ID        string `json:"id"`
	TwitterID string `json:"twitterId"`

	Name        string `json:"name"`
	TwitterName string `json:"twitterName"`
}
