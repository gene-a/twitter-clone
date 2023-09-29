package models

type Tweet struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Tweet  string `json:"tweet"`
}