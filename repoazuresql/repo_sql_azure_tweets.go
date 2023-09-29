package repoazuresql

import (
	"database/sql"
	"twitter_clone/models"
)

func CreateTweet(tweet models.Tweet) error {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "INSERT INTO Tweets (UserID, Tweet) VALUES (@UserID, @Tweet)"
	_, err := db.Exec(query, sql.Named("UserID", tweet.UserID), sql.Named("Tweet", tweet.Tweet))
	if err != nil {
		return err
	}

	return nil
}

func GetAllTweets() ([]models.Tweet, error) {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "SELECT * FROM Tweets"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Tweet)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

func GetTweetFeedByUserId(userId int64, page, pageSize int) ([]models.FeedItem, error) {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	// Calculate the offset based on the page and pageSize
	offset := (page - 1) * pageSize

	// Use OFFSET and FETCH clauses for the query to simulate pagination
	// ID Descending as we need to get the latest tweets from all the users we are subscribed to
	// We also include our own tweets in the feed
	query := "SELECT Tweets.Id AS Id, Tweets.Tweet AS Tweet, Users.Username, Tweets.UserID FROM Tweets  JOIN Followers ON Tweets.UserId = Followers.UserId JOIN Users ON Tweets.UserId = Users.Id WHERE Followers.FollowerUserId = @UserID OR Tweets.UserId = @UserId ORDER BY ID DESC OFFSET @Offset ROWS FETCH NEXT @PageSize ROWS ONLY"

	rows, err := db.Query(query, sql.Named("UserID", userId), sql.Named("Offset", offset), sql.Named("PageSize", pageSize))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feedItems []models.FeedItem
	for rows.Next() {
		var feedItem models.FeedItem
		err := rows.Scan(&feedItem.ID, &feedItem.Tweet, &feedItem.Username, &feedItem.UserID)
		if err != nil {
			return nil, err
		}
		feedItems = append(feedItems, feedItem)
	}

	return feedItems, nil
}

func GetSuggestedTweet(userId int64) (models.FeedItem, error) {
    db := GetAzureSqlDbConnection()
    defer db.Close()

    // Get a RANDOM tweet from someone who the userId does not follow
    query := "SELECT Top 1 Tweets.Id, Tweets.Tweet, Users.Username, Tweets.UserID FROM Users JOIN Tweets ON Users.Id = Tweets.UserId WHERE Users.Id <> @UserId AND Users.Id NOT IN (SELECT Followers.UserId FROM Followers WHERE Followers.FollowerUserId = @UserId) ORDER BY NEWID()"

    row := db.QueryRow(query, sql.Named("UserID", userId))

    var feedItem models.FeedItem
    err := row.Scan(&feedItem.ID, &feedItem.Tweet, &feedItem.Username, &feedItem.UserID)
    if err != nil {
        if err == sql.ErrNoRows {
            // Handle the case when there are no suggested tweets
            return models.FeedItem{}, nil
        }
        // Handle other errors
        return models.FeedItem{}, err
    }

    return feedItem, nil
}

func GetUserTweetsByUserId(userId int64) ([]models.Tweet, error) {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "SELECT * FROM Tweets WHERE UserId=@UserID"
	rows, err := db.Query(query, sql.Named("UserID", userId))

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Tweet)
		if err != nil {
			return nil, err
		}
		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

func DeleteTweet(tweetId int64, userId int64) error {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "DELETE FROM Tweets WHERE ID=@TweetID AND UserId=@UserID"
	_, err := db.Exec(query, sql.Named("TweetID", tweetId), sql.Named("UserID", userId))
	if err != nil {
		return err
	}

	return nil
}
