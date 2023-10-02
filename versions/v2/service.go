package v2

import (
	"database/sql"
	"layer/config"
	"log"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUser registers a new user in the database and returns the user's data.
func RegisterUser(email string, password string) (*User, error) {
	if len(password) < 8 {
		return nil, &ValidationError{
			Field:   "password",
			Message: "Password must be at least 8 characters long",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return nil, err
	}

	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		return nil, err
	}
	defer db.Close()

	// Use a prepared statement
	stmt, err := db.Prepare("INSERT INTO users (email, password) VALUES (?, ?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		return nil, err
	}
	defer stmt.Close()

	userToCreate := CreateUserParams{
		Email:    email,
		Password: string(hashedPassword),
	}

	_, err = stmt.Exec(userToCreate.Email, userToCreate.Password)
	if err != nil {
		log.Printf("Failed to execute statement: %v", err)
		return nil, err
	}

	var userFromDB User
	err = db.QueryRow("SELECT * FROM users WHERE email=?", email).Scan(
		&userFromDB.ID,
		&userFromDB.Email,
		&userFromDB.Password,
		&userFromDB.CreatedAt,
	)
	if err != nil {
		log.Printf("Failed to query database: %v", err)
		return nil, err
	}

	return &userFromDB, nil
}
