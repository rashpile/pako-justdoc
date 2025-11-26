package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rashpile/pako-justdoc/internal/api"
	"github.com/rashpile/pako-justdoc/internal/storage"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "justdoc.db"
	}

	// Initialize storage
	store, err := storage.NewBoltStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		_ = store.Close()
	}()

	// Initialize API
	handler := api.NewHandler(store)
	router := api.NewRouter(handler)

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\nShutting down...")
		_ = store.Close()
		os.Exit(0)
	}()

	// Start server
	fmt.Printf("JustDoc starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
