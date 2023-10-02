package v4

import (
	"fmt"
	"layer/config"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func Run() {
	repo, err := NewRepository(config.DBPath)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	svc := NewService(repo)
	router := CreateRouter(svc)

	fmt.Println("Server listening on", config.ServerAddress)
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
