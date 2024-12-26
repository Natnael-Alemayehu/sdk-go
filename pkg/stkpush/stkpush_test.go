package stkpush

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockClient struct{}

func (m *MockClient) DoRequest(method, endpoint string, body interface{}) ([]byte, error) {
	response := STKPushResponse{
		MerchantRequestID:   "12345",
		CheckoutRequestID:   "67890",
		ResponseCode:        "0",
		ResponseDescription: "Success",
		CustomerMessage:     "Request accepted for processing",
	}
	return json.Marshal(response)
}

func TestInitiateSTKPush(t *testing.T) {
	mockClient := &MockClient{}
	service := NewSTKPushService(mockClient)

	request := &STKPushRequest{
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

	response, err := service.InitiateSTKPush(request)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "12345", response.MerchantRequestID)
	assert.Equal(t, "67890", response.CheckoutRequestID)
	assert.Equal(t, "0", response.ResponseCode)
	assert.Equal(t, "Success", response.ResponseDescription)
	assert.Equal(t, "Request accepted for processing", response.CustomerMessage)
}
