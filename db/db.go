package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

const createTableSQL = `
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

// ResetDB initializes or resets the database to its initial state.
func ResetDB(databasePath string) {
	// Delete the database file if it exists
	os.Remove(databasePath)

	// Create the database file
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
}

const queryUsersSQL = `
SELECT
	id,
	email,
	password,
	created_at
FROM users;
`

// ShowDB displays the contents of the database.
func ShowDB(databasePath string) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	rows, err := db.Query(queryUsersSQL)
	if err != nil {
		log.Fatalf("Failed to query users: %v", err)
	}
	defer rows.Close()

	// Scan the rows into a slice of users
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
		)
		if err != nil {
			log.Fatalf("Failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	if users == nil {
		fmt.Println("No users found")
		return
	}

	// Print the users as JSON lines
	for _, user := range users {
		userJSON, _ := json.Marshal(user)
		fmt.Println(string(userJSON))
	}
}
