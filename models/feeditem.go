package models

type FeedItem struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	UserID   string `json:"user_id"`
	Tweet    string `json:"tweet"`
}