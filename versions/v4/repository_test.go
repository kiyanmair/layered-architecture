package v4

import (
	"layer/config"
	"layer/db"
	"os"
	"testing"
)

func NewTestRepository() (Repository, error) {
	db.ResetDB(config.TestDBPath)
	return NewRepository(config.TestDBPath)
}

func TestUsersRepository(t *testing.T) {
	testCases := []struct {
		name  string
		email string
	}{
		{
			name:  "first user",
			email: "test1@example.com",
		},
		{
			name:  "second user",
			email: "test2@example.com",
		},
	}

	repo, err := NewTestRepository()
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	// Defers are executed in LIFO order
	defer os.Remove(config.TestDBPath)
	defer repo.Close()

	// Keep track of unique IDs returned by the repository
	uniqueIDs := make(map[int]bool)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userToCreate := &CreateUserParams{
				Email:    tc.email,
				Password: "securepass123",
			}

			err = repo.CreateUser(userToCreate)
			if err != nil {
				t.Fatalf("Failed to create user: %v", err)
			}

			err = repo.CreateUser(userToCreate)
			if err == nil {
				t.Fatalf("Repository created user with duplicate email")
			}

			userFromDB, err := repo.GetUserByEmail(userToCreate.Email)
			if err != nil {
				t.Fatalf("Failed to get user: %v", err)
			}

			if userFromDB.Email != userToCreate.Email {
				t.Fatalf("Repository returned unexpected email: got %v want %v", userFromDB.Email, userToCreate.Email)
			}

			uniqueIDs[userFromDB.ID] = true
		})
	}

	if len(uniqueIDs) != len(testCases) {
		t.Fatalf("Repository returned duplicate user IDs")
	}
}
