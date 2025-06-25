package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/NanobyteRuata/go-taskmanager/internal/api"
	"github.com/joho/godotenv"
)

const (
	defaultPort = "8080"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it. Using default values.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	handler := api.NewHandler()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler.Router(),
	}

	// Channel to listen for errors coming from the server
	serverErrors := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		log.Printf("Server listening on port %s", port)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until we receive a signal or an error
	select {
	case err := <-serverErrors:
		log.Fatalf("Error starting server: %v", err)
	case <-shutdown:
		log.Println("Shutting down server...")
		// Gracefully shutdown the server
		if err := server.Close(); err != nil {
			log.Fatalf("Could not stop server gracefully: %v", err)
		}
	}

	fmt.Println("Server stopped")
}
