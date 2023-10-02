package v4

import "time"

// User represents a user in the database.
type User struct {
	ID        int       `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}

// CreateUserParams represents the expected fields when creating a user.
type CreateUserParams struct {
	Email    string `db:"email"`
	Password string `db:"password"`
}
