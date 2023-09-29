package repoazuresql

import (
	"database/sql"
	"twitter_clone/models"
)

// CreateFollower creates a new follower record in the database.
func CreateFollower(follower models.Follower) error {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "INSERT INTO Followers (UserID, FollowerUserID) VALUES (@UserID, @FollowerUserID)"
	_, err := db.Exec(query, sql.Named("UserID", follower.UserID), sql.Named("FollowerUserID", follower.FollowerUserID))
	if err != nil {
		return err
	}
	

	return nil
}

// GetFollowersByUserId retrieves all follower records from the database.
func GetFollowersByUserId(userId int64) ([]models.Follower, error) {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "SELECT * FROM Followers WHERE UserID=@UserID"
	rows, err := db.Query(query, sql.Named("UserID", userId))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []models.Follower
	for rows.Next() {
		var follower models.Follower
		err := rows.Scan(&follower.UserID, &follower.FollowerUserID)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	return followers, nil
}

// DeleteFollower deletes a follower record from the database by user ID and follower user ID.
func DeleteFollowerFromUser(userId, followerUserId int64) error {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "DELETE FROM Followers WHERE UserID=@UserID AND FollowerUserID=@FollowerUserID"
	_, err := db.Exec(query, sql.Named("UserID", userId), sql.Named("FollowerUserID", followerUserId))
	if err != nil {
		return err
	}

	return nil
}
