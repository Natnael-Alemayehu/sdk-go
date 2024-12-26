// cmd/example/main.go

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/client"
	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/config"
)

func main() {
	// Get credentials from environment variables
	// consumerKey := os.Getenv("MPESA_CONSUMER_KEY")
	// consumerSecret := os.Getenv("MPESA_CONSUMER_SECRET")

	// if consumerKey == "" || consumerSecret == "" {
	// 	log.Fatal("MPESA_CONSUMER_KEY and MPESA_CONSUMER_SECRET environment variables are required")
	// }

	// Initialize config
	cfg := &config.Config{
		BaseURL:        "https://apisandbox.safaricom.et",
		ConsumerKey:    "9AevUoujV91YANTnVaokwlt4TOD8H9zxNLsa1I1xwaWsA2Qm",
		ConsumerSecret: "2GwJD3HAyrsBug1PzfXmNUNDfWhgKAOuOJm0sF7Ct4bpYawszSxmUxrRAxTalELp",
		Timeout:        time.Second * 30,
		RetryCount:     3,
		RetryWaitTime:  time.Second * 2,
	}

	// Create client
	client := client.NewClient(cfg)

	// Test authentication
	fmt.Println("Testing authentication...")

	// The client will automatically handle authentication when making a request
	if err := client.GetToken(); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	fmt.Println("Basic Authentication successful!")
}
