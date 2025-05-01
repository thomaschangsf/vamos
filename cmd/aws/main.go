package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/chang/vamos/internal/aws"
)

/*


Usage:
  Option1: Build binary with:
		make build-aws
		./bin/aws [command] [options]
		./bin/aws -region us-west-2 session 1h
		./bin/aws -region us-west-2 buckets s3://commoncrawl/crawl-data/CC-MAIN-2024-10/segments/
  Option2: go run cmd/aws/main.go [command] [options]

Commands:
  session     Get temporary session token and export to environment
  identity    Show AWS identity information
  buckets     List S3 buckets or objects in a bucket
  help        Show help message

Examples:
  go run cmd/aws/main.go identity
  go run cmd/aws/main.go buckets
  go run cmd/aws/main.go session
  go run cmd/aws/main.go help


To run tests:
	go test -v ./internal/aws/...
*/

func main() {
	region := flag.String("region", "us-west-2", "AWS region")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: aws [options] <command>")
		fmt.Println("Commands:")
		fmt.Println("  buckets [s3://bucket/path]  List S3 buckets or objects in a bucket")
		fmt.Println("  identity                    Get AWS identity information")
		fmt.Println("  session [duration]          Get temporary session token")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := aws.NewClient(ctx, *region)
	if err != nil {
		fmt.Printf("Error creating AWS client: %v\n", err)
		os.Exit(1)
	}

	command := flag.Arg(0)
	switch command {
	case "buckets":
		handleBuckets(ctx, client, flag.Arg(1))
	case "identity":
		handleIdentity(ctx, client)
	case "session":
		handleSession(ctx, client, flag.Arg(1))
	case "help":
		printHelp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printHelp()
		os.Exit(1)
	}
}

func handleIdentity(ctx context.Context, client *aws.Client) {
	identity, err := client.GetCallerIdentity(ctx)
	if err != nil {
		fmt.Printf("Error getting identity: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("AWS Identity:")
	fmt.Printf("  Account: %s\n", *identity.Account)
	fmt.Printf("  ARN:     %s\n", *identity.Arn)
	fmt.Printf("  User ID: %s\n", *identity.UserId)
}

func handleBuckets(ctx context.Context, client *aws.Client, location string) {
	if location == "" {
		// List all buckets
		buckets, err := client.ListBuckets(ctx)
		if err != nil {
			fmt.Printf("Error listing buckets: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("S3 Buckets:")
		for _, bucket := range buckets {
			fmt.Printf("  %s\n", *bucket.Name)
		}
	} else {
		// List objects in the specified bucket
		objects, err := client.ListObjects(ctx, location)
		if err != nil {
			fmt.Printf("Error listing objects: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Objects in %s:\n", location)
		for _, obj := range objects {
			fmt.Printf("  %s\n", *obj.Key)
		}
	}
}

func handleSession(ctx context.Context, client *aws.Client, durationStr string) {
	duration := int32(3600) // Default duration: 1 hour
	if durationStr != "" {
		d, err := time.ParseDuration(durationStr)
		if err != nil {
			fmt.Printf("Invalid duration: %v\n", err)
			os.Exit(1)
		}
		duration = int32(d.Seconds())
	}

	credentials, err := client.GetSessionToken(ctx, duration)
	if err != nil {
		fmt.Printf("Error getting session token: %v\n", err)
		os.Exit(1)
	}

	// Export credentials to environment
	os.Setenv("AWS_ACCESS_KEY_ID", *credentials.Credentials.AccessKeyId)
	os.Setenv("AWS_SECRET_ACCESS_KEY", *credentials.Credentials.SecretAccessKey)
	os.Setenv("AWS_SESSION_TOKEN", *credentials.Credentials.SessionToken)

	// Try to set the credentials in the current shell if possible
	cmd := exec.Command("export", fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", *credentials.Credentials.AccessKeyId))
	cmd.Run()
	cmd = exec.Command("export", fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", *credentials.Credentials.SecretAccessKey))
	cmd.Run()
	cmd = exec.Command("export", fmt.Sprintf("AWS_SESSION_TOKEN=%s", *credentials.Credentials.SessionToken))
	cmd.Run()

	fmt.Println("Temporary credentials have been set in the environment.")
	fmt.Printf("Expires at: %s\n", credentials.Credentials.Expiration.Format(time.RFC3339))
}

func printHelp() {
	fmt.Println(`
Usage:
  go run cmd/aws/main.go [command] [options]

Commands:
  identity    Show AWS identity information
  buckets     List S3 buckets or objects in a bucket
  session     Get temporary session token
  help        Show help message

Options:
  -region string
        AWS region (default "us-west-2")

Examples:
  go run cmd/aws/main.go identity
  go run cmd/aws/main.go buckets
  go run cmd/aws/main.go session
  go run cmd/aws/main.go help
`)
}
