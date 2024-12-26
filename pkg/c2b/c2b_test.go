package c2b

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	doRequestFunc func(method, endpoint string, body interface{}) ([]byte, error)
}

func (m *mockClient) DoRequest(method, endpoint string, body interface{}) ([]byte, error) {
	return m.doRequestFunc(method, endpoint, body)
}

func TestProcessPayment(t *testing.T) {
	tests := []struct {
		name           string
		request        *PaymentRequest
		mockResponse   []byte
		mockError      error
		expectedResult *PaymentResponse
		expectedError  error
	}{
		{
			name: "successful payment",
			request: &PaymentRequest{
				RequestRefID: "12345",
				CommandID:    "CustomerPayBillOnline",
				SourceSystem: "USSD",
			},
			mockResponse: json.RawMessage(`{
				"RequestRefID": "12345",
				"ResponseCode": "0",
				"ResponseDesc": "Success",
				"TransactionID": "67890",
				"AdditionalInfo": []
			}`),
			mockError: nil,
			expectedResult: &PaymentResponse{
				RequestRefID:   "12345",
				ResponseCode:   "0",
				ResponseDesc:   "Success",
				TransactionID:  "67890",
				AdditionalInfo: []string{},
			},
			expectedError: nil,
		},
		{
			name: "failed payment",
			request: &PaymentRequest{
				RequestRefID: "12345",
				CommandID:    "CustomerPayBillOnline",
				SourceSystem: "USSD",
			},
			mockResponse:   nil,
			mockError:      errors.New("failed to process C2B payment"),
			expectedResult: nil,
			expectedError:  errors.New("failed to process C2B payment: failed to process C2B payment"),
		},
		{
			name: "invalid response",
			request: &PaymentRequest{
				RequestRefID: "12345",
				CommandID:    "CustomerPayBillOnline",
				SourceSystem: "USSD",
			},
			mockResponse:   json.RawMessage(`invalid response`),
			mockError:      nil,
			expectedResult: nil,
			expectedError:  errors.New("failed to parse C2B payment response: invalid character 'i' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockClient{
				doRequestFunc: func(method, endpoint string, body interface{}) ([]byte, error) {
					return tt.mockResponse, tt.mockError
				},
			}

			service := NewC2BService(mockClient)
			result, err := service.ProcessPayment(tt.request)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResult, result)
		})
	}
}
