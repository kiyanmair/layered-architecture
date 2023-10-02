package v4

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockService struct{}

func (m *MockService) RegisterUser(email, password string) (*User, error) {
	user := &User{
		ID:        1,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
	return user, nil
}

func TestRegisterUserHandler(t *testing.T) {
	mockSvc := &MockService{}
	handler := RegisterUserHandler(mockSvc)

	body := map[string]string{
		"email":    "test@example.com",
		"password": "securepass123",
	}
	bodyBytes, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "/api/users/register", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := map[string]interface{}{
		"id":    float64(1), // JSON numbers are parsed as floats
		"email": "test@example.com",
	}

	var actual map[string]interface{}
	err = json.NewDecoder(recorder.Body).Decode(&actual)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if actual["id"] != expected["id"] {
		t.Errorf("Handler returned unexpected ID: got %v want %v", actual["id"], expected["id"])
	}

	if actual["email"] != expected["email"] {
		t.Errorf("Handler returned unexpected email: got %v want %v", actual["email"], expected["email"])
	}
}

func TestRegisterUserHandler_InvalidInput(t *testing.T) {
	testCases := []struct {
		name string
		body map[string]string
	}{
		{
			name: "invalid body",
			body: map[string]string{
				"test": "test",
			},
		},
		{
			name: "missing email",
			body: map[string]string{
				"password": "securepass123",
			},
		},
		{
			name: "missing password",
			body: map[string]string{
				"email": "test@example.com",
			},
		},
		{
			name: "invalid email",
			body: map[string]string{
				"email":    "test",
				"password": "securepass123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := &MockService{}
			handler := RegisterUserHandler(mockSvc)

			bodyBytes, _ := json.Marshal(tc.body)

			req, err := http.NewRequest("POST", "/api/users/register", bytes.NewReader(bodyBytes))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			status := recorder.Code
			if status != http.StatusBadRequest {
				t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
			}
		})
	}
}
