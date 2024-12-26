package c2b

import (
	"encoding/json"
	"fmt"

	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/models"
)

type C2BService struct {
	client interface {
		DoRequest(method, endpoint string, body interface{}) ([]byte, error)
	}
}

func NewC2BService(client interface {
	DoRequest(method, endpoint string, body interface{}) ([]byte, error)
}) *C2BService {
	return &C2BService{
		client: client,
	}
}

// RegisterURLRequest represents the request to register C2B callback URLs
type RegisterURLRequest struct {
	ShortCode       string `json:"ShortCode"`
	ResponseType    string `json:"ResponseType"`
	CommandID       string `json:"CommandID"`
	ConfirmationURL string `json:"ConfirmationURL"`
	ValidationURL   string `json:"ValidationURL"`
}

// RegisterURLResponse represents the response from URL registration
type RegisterURLResponse struct {
	Header struct {
		ResponseCode    int    `json:"responseCode"`
		ResponseMessage string `json:"responseMessage"`
		CustomerMessage string `json:"customerMessage"`
		Timestamp       string `json:"timestamp"`
	} `json:"header"`
}

// PaymentRequest represents a C2B payment request
type PaymentRequest struct {
	RequestRefID     string               `json:"RequestRefID"`
	CommandID        string               `json:"CommandID"`
	Remark           string               `json:"Remark"`
	ChannelSessionID string               `json:"ChannelSessionID"`
	SourceSystem     string               `json:"SourceSystem"`
	Timestamp        string               `json:"Timestamp"`
	Parameters       []models.Parameter   `json:"Parameters"`
	ReferenceData    []models.Reference   `json:"ReferenceData"`
	Initiator        models.Initiator     `json:"Initiator"`
	PrimaryParty     models.Party         `json:"PrimaryParty"`
	ReceiverParty    models.ReceiverParty `json:"ReceiverParty"`
}

// PaymentResponse represents a C2B payment response
type PaymentResponse struct {
	RequestRefID   string   `json:"RequestRefID"`
	ResponseCode   string   `json:"ResponseCode"`
	ResponseDesc   string   `json:"ResponseDesc"`
	TransactionID  string   `json:"TransactionID"`
	AdditionalInfo []string `json:"AdditionalInfo"`
}

// RegisterURL registers the confirmation and validation URLs
func (s *C2BService) RegisterURL(req *RegisterURLRequest) (*RegisterURLResponse, error) {
	if req.CommandID == "" {
		req.CommandID = "RegisterURL"
	}
	if req.ResponseType == "" {
		req.ResponseType = "Completed"
	}

	endpoint := "/v1/c2b-register-url/register"
	resp, err := s.client.DoRequest("POST", endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to register URLs: %w", err)
	}

	var urlResp RegisterURLResponse
	if err := json.Unmarshal(resp, &urlResp); err != nil {
		return nil, fmt.Errorf("failed to parse URL registration response: %w", err)
	}

	return &urlResp, nil
}

// ProcessPayment processes a C2B payment
func (s *C2BService) ProcessPayment(req *PaymentRequest) (*PaymentResponse, error) {
	if req.CommandID == "" {
		req.CommandID = "CustomerPayBillOnline"
	}
	if req.SourceSystem == "" {
		req.SourceSystem = "USSD"
	}

	endpoint := "/v1/c2b/payments"
	resp, err := s.client.DoRequest("POST", endpoint, req)
	if err != nil {
		return nil, fmt.Errorf("failed to process C2B payment: %w", err)
	}

	var payResp PaymentResponse
	if err := json.Unmarshal(resp, &payResp); err != nil {
		return nil, fmt.Errorf("failed to parse C2B payment response: %w", err)
	}

	return &payResp, nil
}
