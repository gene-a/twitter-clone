package repoazuresql

import (
	"database/sql"
	"twitter_clone/models"
)

// CreateUser creates a new user record in the database.
func CreateUser(user models.User) error {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "INSERT INTO Users (Email, Username, FirstName, LastName, PasswordHash, DateOfBirth) VALUES (@Email, @Username, @FirstName, @LastName, @PasswordHash, @DateOfBirth)"
	_, err := db.Exec(query, sql.Named("Email", user.Email), sql.Named("Username", user.Username), sql.Named("FirstName", user.FirstName), sql.Named("LastName", user.LastName), sql.Named("PasswordHash", user.PasswordHash), sql.Named("DateOfBirth", user.DateOfBirth))
	if err != nil {
		return err
	}

	return nil
}

// GetAllUsers retrieves all user records from the database.
func GetAllUsers() ([]models.User, error) {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "SELECT ID, Email, Username, FirstName, LastName, DateOfBirth FROM Users"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		// Skipping over returning the password
		err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.DateOfBirth)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserById retrieves a user record from the database by ID.
func GetUserById(userId int64) (models.User, error) {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "SELECT ID, Email, Username, FirstName, LastName, DateOfBirth FROM [dbo].[Users] WHERE ID = @UserID"  // @UserID is a named parameter
	row := db.QueryRow(query, sql.Named("UserID", userId))

	var user models.User
	// Skipping over returning the password
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.DateOfBirth)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// GetUserById retrieves a user record from the database by ID.
func GetUserByUsernameAndPassword(username string, password string) (models.User, error) {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "SELECT ID, Email, Username, FirstName, LastName, DateOfBirth FROM [dbo].[Users] WHERE Username = @Username AND PasswordHash = @Passwordhash"  
	row := db.QueryRow(query, sql.Named("Username", username), sql.Named("Passwordhash", password))

	var user models.User
	// Skipping over returning the password
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.FirstName, &user.LastName, &user.DateOfBirth)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

// UpdateUser updates an existing user record in the database.
func UpdateUser(user models.User) error {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "UPDATE Users SET Email=@Email, Username=@Username, FirstName=@FirstName, LastName=@LastName, PasswordHash=@PasswordHash, DateOfBirth=@DateOfBirth WHERE ID=@ID"
	_, err := db.Exec(query, sql.Named("Email", user.Email), sql.Named("Username", user.Username), sql.Named("FirstName", user.FirstName), sql.Named("LastName", user.LastName), sql.Named("PasswordHash", user.PasswordHash), sql.Named("DateOfBirth", user.DateOfBirth), sql.Named("ID", user.ID))
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user record from the database by ID.
func DeleteUser(userID int64) error {
	db := GetAzureSqlDbConnection()
	defer db.Close()

	query := "DELETE FROM Users WHERE ID=@ID"
	_, err := db.Exec(query, sql.Named("ID", userID))
	if err != nil {
		return err
	}

	return nil
}
