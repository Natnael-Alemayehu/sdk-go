package stkpush

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/models"
)

type STKPushService struct {
	client interface {
		DoRequest(method, endpoint string, body interface{}) ([]byte, error)
	}
}

func NewSTKPushService(client interface {
	DoRequest(method, endpoint string, body interface{}) ([]byte, error)
}) *STKPushService {
	return &STKPushService{
		client: client,
	}
}

type STKPushRequest struct {
	MerchantRequestID string                 `json:"MerchantRequestID"`
	BusinessShortCode string                 `json:"BusinessShortCode"`
	Password          string                 `json:"Password"`
	Timestamp         string                 `json:"Timestamp"`
	TransactionType   string                 `json:"TransactionType"`
	Amount            string                 `json:"Amount"`
	PartyA            string                 `json:"PartyA"`
	PartyB            string                 `json:"PartyB"`
	PhoneNumber       string                 `json:"PhoneNumber"`
	TransactionDesc   string                 `json:"TransactionDesc"`
	CallBackURL       string                 `json:"CallBackURL"`
	AccountReference  string                 `json:"AccountReference"`
	ReferenceData     []models.ReferenceItem `json:"ReferenceData,omitempty"`
}

type STKPushResponse struct {
	MerchantRequestID   string `json:"MerchantRequestID"`
	CheckoutRequestID   string `json:"CheckoutRequestID"`
	ResponseCode        string `json:"ResponseCode"`
	ResponseDescription string `json:"ResponseDescription"`
	CustomerMessage     string `json:"CustomerMessage"`
}

func (s *STKPushService) InitiateSTKPush(req *STKPushRequest) (*STKPushResponse, error) {
	if req.Timestamp == "" {
		req.Timestamp = time.Now().Format("20060102150405") // Format: YYYYMMDDHHMMSS
	}

	if req.TransactionType == "" {
		req.TransactionType = "CustomerPayBillOnline"
	}

	endpoint := "/mpesa/stkpush/v3/processrequest"
	resp, err := s.client.DoRequest("POST", endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("STK push request failed: %w", err)
	}

	var stkResp STKPushResponse
	if err := json.Unmarshal(resp, &stkResp); err != nil {
		return nil, fmt.Errorf("failed to parse STK push response: %w", err)
	}

	return &stkResp, nil
}
