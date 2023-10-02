package v4

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(email string, password string) (*User, error)
}

type UserService struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &UserService{
		repo: repo,
	}
}

// RegisterUser registers a new user and returns the user's data.
func (s *UserService) RegisterUser(email string, password string) (*User, error) {
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

	userToCreate := &CreateUserParams{
		Email:    email,
		Password: string(hashedPassword),
	}

	err = s.repo.CreateUser(userToCreate)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return nil, err
	}

	userFromDB, err := s.repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("Failed to get user: %v", err)
		return nil, err
	}

	return userFromDB, nil
}
