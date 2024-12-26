package auth

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/models"
)

type AuthService struct {
	client interface {
		DoRequest(method, endpoint string, body interface{}) ([]byte, error)
	}
	mu       sync.RWMutex
	token    string
	tokenExp time.Time
}

func NewAuthService(client interface {
	DoRequest(method, endpoint string, body interface{}) ([]byte, error)
}) *AuthService {
	return &AuthService{
		client: client,
	}
}

func (s *AuthService) GetToken() (string, error) {
	s.mu.RLock()

	if s.token != "" && time.Now().Before(s.tokenExp) {
		fmt.Printf("The token is still valid: %s\n", s.token)
		defer s.mu.RUnlock()
		return s.token, nil
	}
	s.mu.RUnlock()
	return s.refreshToken()
}

func (s *AuthService) refreshToken() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check token expiration after acquiring lock
	if s.token != "" && time.Now().Before(s.tokenExp) {
		return s.token, nil
	}

	endpoint := "/v1/token/generate?grant_type=client_credentials"
	resp, err := s.client.DoRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("failed to refresh token: %w", err)
	}

	var authResp models.AuthResponse
	if err := json.Unmarshal(resp, &authResp); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %w", err)
	}

	s.token = authResp.AccessToken
	// Convert expires_in to duration and subtract a minute for safety
	expiresIn := 3600 * time.Second // Default to 1 hour if parsing fails
	s.tokenExp = time.Now().Add(expiresIn - time.Minute)

	return s.token, nil
}
