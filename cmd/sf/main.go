package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/thomaschangsf/vamos/internal/aws"
	"github.com/thomaschangsf/vamos/internal/llm"
	"github.com/thomaschangsf/vamos/internal/web"
	"github.com/thomaschangsf/vamos/pkg/config"
)

func main() {
	// Load configuration
	cfg := config.NewConfig()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize AWS client
	awsClient, err := aws.NewClient(ctx, cfg.AWSRegion)
	if err != nil {
		log.Fatalf("Failed to create AWS client: %v", err)
	}

	// Initialize LLM client
	llmClient := llm.NewClient(cfg.LLMAPIKey, cfg.LLMModelName)

	// Initialize web server
	webServer := web.NewServer(cfg.WebPort)

	// Start web server in a goroutine
	go func() {
		if err := webServer.Start(ctx); err != nil {
			log.Printf("Web server error: %v", err)
		}
	}()

	// Example: List S3 buckets
	buckets, err := awsClient.ListBuckets(ctx)
	if err != nil {
		log.Printf("Failed to list buckets: %v", err)
	} else {
		fmt.Println("S3 Buckets:")
		for _, bucket := range buckets {
			fmt.Printf("- %s\n", bucket)
		}
	}

	// Example: Generate text with LLM
	response, err := llmClient.GenerateText(ctx, "What is the capital of France?")
	if err != nil {
		log.Printf("Failed to generate text: %v", err)
	} else {
		fmt.Printf("\nLLM Response: %s\n", response)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Cancel context to stop all services
	cancel()
}
