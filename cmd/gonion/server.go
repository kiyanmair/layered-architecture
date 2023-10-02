package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"layer/config"
	"net/http"
	"time"
)

// ServerRunning checks if the server is running and prints a message if not.
func ServerRunning() bool {
	if !PingServer() {
		fmt.Println("Server is not running. Please use `gonion run` and try again.")
		return false
	}
	return true
}

// PingServer sends a simple request to the server to check if it's running.
func PingServer() bool {
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := client.Get("http://" + config.ServerAddress + "/api/ping")
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

// RequestData represents a generic structure for sending data in requests.
type RequestData map[string]interface{}

// SendRequest sends an HTTP request to the server.
func SendRequest(endpoint string, method string, data RequestData) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %v", err)
	}

	req, err := http.NewRequest(
		method,
		fmt.Sprintf("http://%s%s", config.ServerAddress, endpoint),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

// HandleResponse checks the response status and reads the body.
func HandleResponse(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
