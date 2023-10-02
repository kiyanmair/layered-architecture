package v3

// RegisterUserRequest represents the expected format for a user registration request.
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterUserResponse represents the response format after successfully registering a user.
type RegisterUserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
