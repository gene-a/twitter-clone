package models

type Follower struct {
	UserID         int64 `json:"user_id"`
	FollowerUserID int64 `json:"follower_user_id"`
}