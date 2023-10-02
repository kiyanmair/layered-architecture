package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Registers a new user by sending an HTTP request to the server.
func RegisterUser(email, password string) error {
	requestData := RequestData{
		"email":    email,
		"password": password,
	}

	resp, err := SendRequest("/api/users/register", "POST", requestData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = HandleResponse(resp)
	return err
}

// PromptPassword prompts the user for a password and returns it.
func PromptPassword() (string, error) {
	fmt.Print("Password: ")

	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to read password: %v", err)
	}
	fmt.Println() // newline after password

	password := strings.TrimSpace(string(bytePassword))
	return password, nil
}
