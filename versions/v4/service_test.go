package v4

import (
	"errors"
	"testing"
	"time"
)

type MockRepository struct {
	Users map[string]*User
}

func (m *MockRepository) CreateUser(user *CreateUserParams) error {
	mockUser := &User{
		ID:        1,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
	}

	if m.Users == nil {
		m.Users = make(map[string]*User)
	}
	m.Users[mockUser.Email] = mockUser

	return nil
}

func (m *MockRepository) GetUserByEmail(email string) (*User, error) {
	mockUser, exists := m.Users[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return mockUser, nil
}

func (m *MockRepository) Close() error {
	return nil
}

func TestRegisterUser(t *testing.T) {
	mockRepo := &MockRepository{}
	svc := NewService(mockRepo)

	email := "test@example.com"
	password := "securepass123"

	user, err := svc.RegisterUser(email, password)
	if err != nil {
		t.Fatalf("Failed to register user: %v", err)
	}

	if user.Email != email {
		t.Errorf("Service returned unexpected email: got %v want %v", user.Email, email)
	}

	if user.Password == password {
		t.Errorf("Service returned unhashed password")
	}
}

func TestRegisterUser_InvalidInput(t *testing.T) {
	testCases := []struct {
		name     string
		field    string
		email    string
		password string
	}{
		{
			name:     "short password",
			field:    "password",
			email:    "test@example.com",
			password: "short",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockRepository{}
			svc := NewService(mockRepo)

			_, err := svc.RegisterUser(tc.email, tc.password)
			var validationErr *ValidationError
			if errors.As(err, &validationErr) {
				if validationErr.Field != tc.field {
					t.Errorf("Service returned unexpected validation error field: got %v want %v", validationErr.Field, tc.field)
				}
			} else {
				t.Errorf("Service did not return validation error")
			}
		})
	}
}
