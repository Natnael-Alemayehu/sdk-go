package main

import (
	"log"
	"time"

	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/c2b"
	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/client"
	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/config"
	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/models"
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

	// Initialize C2B service
	c2bService := c2b.NewC2BService(client)

	// Register URLs
	urlReq := &c2b.RegisterURLRequest{
		ShortCode:       "101010",
		ConfirmationURL: "http://mydomain.com/c2b/confirmation",
		ValidationURL:   "http://mydomain.com/c2b/validation",
	}

	urlResp, err := c2bService.RegisterURL(urlReq)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("URL Registration Response: %+v\n", urlResp)

	// Process C2B Payment
	payReq := &c2b.PaymentRequest{
		RequestRefID:     "test-ref-id",
		CommandID:        "CustomerPayBillOnline",
		Remark:           "Test Payment",
		ChannelSessionID: "10100000037656400042",
		SourceSystem:     "USSD",
		Timestamp:        time.Now().Format("2006-01-02T15:04:05.000-07:00"),
		Parameters: []models.Parameter{
			{
				Key:   "Amount",
				Value: "500",
			},
			{
				Key:   "AccountReference",
				Value: "TEST",
			},
		},
		Initiator: models.Initiator{
			IdentifierType:     1,
			Identifier:         "251799100026",
			SecurityCredential: "your-security-credential",
			SecretKey:          "your-secret-key",
		},
		PrimaryParty: models.Party{
			IdentifierType: 1,
			Identifier:     "251799100026",
		},
		ReceiverParty: models.ReceiverParty{
			IdentifierType: 4,
			Identifier:     "370360",
			ShortCode:      "370360",
		},
	}

	payResp, err := c2bService.ProcessPayment(payReq)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Payment Response: %+v\n", payResp)
}
