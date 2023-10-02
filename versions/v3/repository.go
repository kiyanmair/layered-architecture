package v3

import (
	"database/sql"
	"fmt"
)

// Repository represents an interface to the database.
type Repository struct {
	db *sql.DB
}

func NewRepository(databasePath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	repo := &Repository{
		db: db,
	}
	return repo, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

const createUserSQL = `
INSERT INTO users (email, password)
VALUES (?, ?);
`

// CreateUser creates a new user in the database.
func (r *Repository) CreateUser(user *CreateUserParams) error {
	// Use a prepared statement
	stmt, err := r.db.Prepare(createUserSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

const getUserByEmailSQL = `
SELECT
	id,
	email,
	password,
	created_at
FROM users
WHERE email=?;
`

// GetUserByEmail returns a user from the database by email.
func (r *Repository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow(getUserByEmailSQL, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	return &user, nil
}
