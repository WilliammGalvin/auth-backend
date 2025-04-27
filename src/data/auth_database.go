package data

import (
	"backend/models"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

// AddUser adds a user to the local database
func AddUser(userData *models.NewUser) error {
	cmdSql := "INSERT INTO users (email, password, display_name) VALUES (?, ?, ?)"
	_, err := database.Exec(cmdSql, userData.Email, userData.Password, userData.DisplayName)

	if err != nil {
		return fmt.Errorf("error adding user with email %s: %w", userData.Email, err)
	}

	return nil
}

// GetUserByEmail fetches user from local database by email
func GetUserByEmail(email string) (*models.User, error) {
	cmdSql := "SELECT * FROM users WHERE email=?"
	rows, err := database.Query(cmdSql, email)

	if err != nil {
		return nil, fmt.Errorf("error fetching user with email %s: %w", email, err)
	}

	defer func(rows *sql.Rows) {
		if err = rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	if !rows.Next() {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	var user models.User
	err = rows.Scan(&user.Id, &user.Email, &user.Password, &user.DisplayName, &user.ProfileImgx64)

	if err != nil {
		return nil, fmt.Errorf("error scanning user with email %s: %w", email, err)
	}

	return &user, nil
}
