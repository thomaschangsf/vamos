package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

func main() {
	fmt.Println("Making HTTP requests...")

	// Using built-in net/http
	fmt.Println("\nUsing net/http:")
	makeHttpRequest()

	// Using resty
	fmt.Println("\nUsing resty:")
	makeRestyRequest()
}

func makeHttpRequest() {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Response status: %s\n", resp.Status)
	fmt.Printf("Response body: %s\n", string(body))
}

func makeRestyRequest() {
	client := resty.New()
	client.SetTimeout(5 * time.Second)

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		Get("https://httpbin.org/get")

	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}

	fmt.Printf("Response status: %s\n", resp.Status())
	fmt.Printf("Response body: %s\n", string(resp.Body()))
}
