package services

import (
	"errors"
	"twitter_clone/models"
	"twitter_clone/repoazuresql"
)


func doesUserExist(userId int64) (models.User, error){
	return repoazuresql.GetUserById(userId)
}

func VerifyUser(userId int64) (error){
	_, err:= doesUserExist(userId)

	if (err != nil){
		return errors.New("user not found")
	}
	return nil
}