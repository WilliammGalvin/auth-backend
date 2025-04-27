package data

import (
	"database/sql"
	"log"
)

var database *sql.DB

// InitDB initializes the local database
func InitDB() {
	db, err := sql.Open("sqlite3", "../data/database.db")

	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	database = db
}

// CloseDB runs appropriate closing logic
func CloseDB() error {
	err := database.Close()

	if err != nil {
		return err
	}

	return nil
}
