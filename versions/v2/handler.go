package v2

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
)

// Ping provides a simple health check endpoint.
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// RegisterUserHandler handles requests to register a new user.
func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	_, err = mail.ParseAddress(req.Email)
	if err != nil {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	user, err := RegisterUser(req.Email, req.Password)
	if err != nil {
		// Handle validation errors
		var validationErr *ValidationError
		if errors.As(err, &validationErr) {
			http.Error(w, validationErr.Error(), http.StatusBadRequest)
			return
		}
		// Handle other errors
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	resp := RegisterUserResponse{
		ID:    user.ID,
		Email: user.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
