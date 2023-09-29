package services

import (
	"twitter_clone/models"
	"twitter_clone/repoazuresql"
)

func CreateUser(user models.User) error {
	return repoazuresql.CreateUser(user)
}

func GetUserById(userId int64) (models.User, error) {
	return repoazuresql.GetUserById(userId)
}

func CreateFollower(follower models.Follower) error {
	// input validation
	// valid user only
	VerifyUser(follower.UserID)
	return repoazuresql.CreateFollower(follower)
}

func GetFollowersByUserId(userId int64) ([]models.Follower, error) {
	// input validation
	// valid user only
	VerifyUser(userId)
	return repoazuresql.GetFollowersByUserId(userId)
}

func DeleteFollowerFromUser(userId, followerUserID int64) error {
	// input validation
	// valid user only
	VerifyUser(userId)
	return repoazuresql.DeleteFollowerFromUser(userId, followerUserID)
}

func GetUserByUsernameAndPassword(username string, password string) (models.User, error) {
	return repoazuresql.GetUserByUsernameAndPassword(username, password)
}