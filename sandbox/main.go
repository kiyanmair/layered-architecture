package sandbox

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"layer/config"
	"log"
	"net/http"
	"net/mail"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func Run() {
	http.HandleFunc("/api/ping", Ping)
	http.HandleFunc("/api/users/register", RegisterUser)
	fmt.Println("Server listening on", config.ServerAddress)
	log.Fatal(http.ListenAndServe(config.ServerAddress, nil))
}

// Ping provides a simple health check endpoint.
func Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("pong"))
}

// RegisterUser handles user registration requests.
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email, emailOk := requestData["email"]
	password, passOk := requestData["password"]
	if !emailOk || !passOk {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	if len(password) < 8 {
		http.Error(w, "Password must be at least 8 characters long", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite3", config.DBPath)
	if err != nil {
		log.Printf("Failed to open database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Use a prepared statement
	stmt, err := db.Prepare("INSERT INTO users (email, password) VALUES (?, ?)")
	if err != nil {
		log.Printf("Failed to prepare statement: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, string(hashedPassword))
	if err != nil {
		log.Printf("Failed to execute statement: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var userID int
	var userEmail string
	err = db.QueryRow("SELECT id, email FROM users WHERE email=$1", email).
		Scan(&userID, &userEmail)
	if err != nil {
		log.Printf("Failed to query database: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]interface{}{
			"id":    userID,
			"email": userEmail,
		},
	)
}
