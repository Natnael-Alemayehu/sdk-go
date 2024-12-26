package auth

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockClient struct{}

func (m *MockClient) DoRequest(method, endpoint string, body interface{}) ([]byte, error) {
	response := map[string]interface{}{
		"access_token": "mock_token",
		"token_type":   "Bearer",
		"expires_in":   "3600",
	}
	return json.Marshal(response)
}

func TestGetToken(t *testing.T) {
	mockClient := &MockClient{}
	authService := NewAuthService(mockClient)

	token, err := authService.GetToken()
	assert.NoError(t, err)
	assert.Equal(t, "mock_token", token)

	// Simulate token expiration
	authService.tokenExp = time.Now().Add(-time.Minute)

	token, err = authService.GetToken()
	assert.NoError(t, err)
	assert.Equal(t, "mock_token", token)
}
