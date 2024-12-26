package main

import (
	"fmt"
	"log"

	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/client"
	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/config"
	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/stkpush"
)

func main() {
	// Initialize config
	cfg, err := config.NewConfig(
		"your-consumer-key",
		"your-consumer-secret",
		config.Sandbox,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create client
	client := client.NewClient(cfg)

	// Initialize STK Push service
	stkService := stkpush.NewSTKPushService(client)

	// Create STK Push request
	request := &stkpush.STKPushRequest{
		BusinessShortCode: "554433",
		Password:          "123",
		Amount:            "10.00",
		PartyA:            "251700404789",
		PartyB:            "554433",
		PhoneNumber:       "251700404789",
		TransactionDesc:   "Test Payment",
		CallBackURL:       "https://your-callback-url.com",
		AccountReference:  "TEST",
	}

	// Initiate STK Push
	response, err := stkService.InitiateSTKPush(request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("STK Push Response: %+v\n", response)
}
