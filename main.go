package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thomaschangsf/vamos/internal/web"
)

func main() {
	// Create a new server on port 8080
	server := web.NewServer(8080)

	// Create a context that we'll use to gracefully shutdown the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	// Start the server
	fmt.Println("Starting server on port 8080...")
	if err := server.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
