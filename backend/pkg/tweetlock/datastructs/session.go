package datastructs

type Session struct {
	SessionID string `dynamodbav:"session_id"`
	UserID    string `dynamodbav:"user_id"`
	ExpiresAt int64  `dynamodbav:"expires_at"`
}
