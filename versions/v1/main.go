package v1

import (
	"fmt"
	"layer/config"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func Run() {
	router := CreateRouter()
	fmt.Println("Server listening on", config.ServerAddress)
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}
