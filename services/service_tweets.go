package services

import (
	"database/sql"
	"errors"
	"twitter_clone/models"
	"twitter_clone/repoazuresql"
)

func GetTweetsByUserId(userId int64) ([]models.Tweet, error) {
	return repoazuresql.GetUserTweetsByUserId(userId)
}

func GetAllTweets(userId int64) ([]models.Tweet, error) {
	return repoazuresql.GetAllTweets()
}

func CreateTweet(tweet models.Tweet) error {
	// input validation
	// valid user only
	err:= VerifyUser(tweet.UserID)

	// tweet length at 280 or less
	if (len(tweet.Tweet) > 280 || err != nil) {
		return errors.New("failed to create tweet")
	}

	return repoazuresql.CreateTweet(tweet)
}

func GetTweetFeedByUserId(userId int64, page, pageSize int) ([]models.FeedItem, error) {
	// input validation
	// valid user only
	VerifyUser(userId)
	return repoazuresql.GetTweetFeedByUserId(userId, page, pageSize)
}

func GetSuggestedTweet(userId int64) (models.FeedItem, error) {
	// input validation
	// valid user only
	VerifyUser(userId)

	// Call the repository function
	feedItem, err := repoazuresql.GetSuggestedTweet(userId)

	if err != nil {
		if err == sql.ErrNoRows {
			// Handle the case when there are no suggested tweets
			// You can return a custom error message or an indicator here.
			return models.FeedItem{}, errors.New("sql.ErrNoRows")
		}
		// Handle other errors
		return models.FeedItem{}, err
	}

	return feedItem, nil
}

func DeleteTweet(tweetId int64, userId int64) error {
	// input validation
	// valid user only
	VerifyUser(userId)
	return repoazuresql.DeleteTweet(tweetId, userId)
}

