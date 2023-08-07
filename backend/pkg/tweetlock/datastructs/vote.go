package datastructs

type Vote struct {
	TargetID string `db:"target_id"`
	UserID   string `db:"user_id"`
}

type VoteCount struct {
	TwitterID string `db:"twitter_id"`
	Votes     uint
	HasVoted  bool `db:"has_voted"`
}
