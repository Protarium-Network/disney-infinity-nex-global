// Command http is a small companion HTTP server (health check) that can run
// next to the NEX servers.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := InitDatabase(); err != nil {
		log.Fatal("Database connection failed: ", err)
	}
	defer CloseDatabase()

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8081"
	}

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           SetupRoutes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("HTTP server listening on :%s", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("HTTP server error: ", err)
	}
}
