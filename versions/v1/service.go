package v1

import (
	"database/sql"
	"layer/config"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUser registers a new user in the database and returns the user's data.
func RegisterUser(email string, password string) (int, string, error) {
	if len(password) < 8 {
		return 0, "", &ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return 0, "", err
	}

	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		return 0, "", err
	}
	defer db.Close()

	// Use a prepared statement
	stmt, err := db.Prepare("INSERT INTO users (email, password) VALUES (?, ?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return 0, "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, hashedPassword)
	if err != nil {
		log.Printf("Failed to execute statement: %v", err)
		return 0, "", err
	}

	var userID int
	var userEmail string
	err = db.QueryRow("SELECT id, email FROM users WHERE email=?", email).
		Scan(&userID, &userEmail)
	if err != nil {
		log.Printf("Failed to query database: %v", err)
		return 0, "", err
	}

	return userID, userEmail, nil
}
