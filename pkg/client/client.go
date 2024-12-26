package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/config"
	"github.com/natnael-alemayehu/mpesa-sdk-go/pkg/models"
)

// Client represents the M-PESA API client
type Client struct {
	config     *config.Config
	httpClient *http.Client
	token      string
	tokenExp   time.Time
}

// NewClient creates a new M-PESA API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// GetToken authenticates with the M-PESA API and gets an access token
func (c *Client) GetToken() error {
	// Check if current token is still valid
	if c.token != "" && time.Now().Before(c.tokenExp) {
		return nil
	}

	// Create basic auth string
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s",
		c.config.ConsumerKey, c.config.ConsumerSecret)))

	// Create URL with query parameters
	baseURL := c.config.BaseURL + "/v1/token/generate"
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("error parsing URL: %w", err)
	}

	q := u.Query()
	q.Add("grant_type", "client_credentials")
	u.RawQuery = q.Encode()

	// Create request
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("error creating auth request: %w", err)
	}

	// Set basic auth header
	req.Header.Set("Authorization", "Basic "+auth)

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making auth request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading auth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Handle specific error codes
		var errorResp struct {
			ResultCode string `json:"resultCode"`
			ResultDesc string `json:"resultDesc"`
		}
		if err := json.Unmarshal(respBody, &errorResp); err != nil {
			return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
		}

		switch errorResp.ResultCode {
		case "999991":
			return fmt.Errorf("invalid client ID: %s", errorResp.ResultDesc)
		case "999996":
			return fmt.Errorf("invalid authentication type: %s", errorResp.ResultDesc)
		case "999997":
			return fmt.Errorf("invalid authorization header: %s", errorResp.ResultDesc)
		case "999998":
			return fmt.Errorf("invalid grant type: %s", errorResp.ResultDesc)
		default:
			return fmt.Errorf("authentication error: %s - %s", errorResp.ResultCode, errorResp.ResultDesc)
		}
	}

	// Parse successful response
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   string `json:"expires_in"`
	}

	if err := json.Unmarshal(respBody, &tokenResp); err != nil {
		return fmt.Errorf("error parsing auth response: %w", err)
	}

	// Store token and expiration
	c.token = tokenResp.AccessToken
	expiresIn := 3599 // default to 1 hour - 1 second if parsing fails
	if tokenResp.ExpiresIn != "" {
		if exp, err := time.ParseDuration(tokenResp.ExpiresIn + "s"); err == nil {
			expiresIn = int(exp.Seconds())
		}
	}
	c.tokenExp = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return nil
}

// DoRequest performs an HTTP request with authentication and retries
func (c *Client) DoRequest(method, endpoint string, body interface{}) ([]byte, error) {
	// Get/refresh token if needed
	if err := c.GetToken(); err != nil {
		return nil, fmt.Errorf("error getting access token: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequest(method, c.config.BaseURL+endpoint, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set common headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	// Implement retry logic
	var lastErr error
	for i := 0; i <= c.config.RetryCount; i++ {
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			time.Sleep(c.config.RetryWaitTime)
			continue
		}
		defer resp.Body.Close()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return respBody, nil
		}

		// Handle error responses
		var errorResp models.CommonResponse
		if err := json.Unmarshal(respBody, &errorResp); err != nil {
			lastErr = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
		} else {
			lastErr = fmt.Errorf("API error: %s - %s", errorResp.ErrorCode, errorResp.ErrorMessage)
		}

		if i < c.config.RetryCount {
			time.Sleep(c.config.RetryWaitTime)
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", c.config.RetryCount, lastErr)
}
