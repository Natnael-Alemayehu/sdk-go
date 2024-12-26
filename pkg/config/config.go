package config

import (
	"fmt"
	"time"
)

// Environment represents the M-PESA API environment
type Environment string

const (
	Production Environment = "production"
	Sandbox    Environment = "sandbox"
)

// Config holds the SDK configuration
type Config struct {
	ConsumerKey    string
	ConsumerSecret string
	Environment    Environment
	BaseURL        string
	Timeout        time.Duration
	RetryCount     int
	RetryWaitTime  time.Duration
}

// ConfigOption defines a function type for configuration options
type ConfigOption func(*Config)

// NewConfig creates a new configuration with the given consumer key and secret
func NewConfig(consumerKey, consumerSecret string, env Environment, options ...ConfigOption) (*Config, error) {
	if consumerKey == "" || consumerSecret == "" {
		return nil, fmt.Errorf("consumer key and secret are required")
	}

	baseURL := "https://apisandbox.safaricom.et"
	if env == Production {
		baseURL = "https://api.safaricom.et"
	}

	cfg := &Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		Environment:    env,
		BaseURL:        baseURL,
		Timeout:        time.Second * 5,
		RetryCount:     2,
		RetryWaitTime:  time.Second * 5,
	}

	// Apply options
	for _, option := range options {
		option(cfg)
	}

	return cfg, nil
}

// WithTimeout sets the client timeout
func WithTimeout(timeout time.Duration) ConfigOption {
	return func(c *Config) {
		c.Timeout = timeout
	}
}

// WithRetry sets retry configuration
func WithRetry(count int, waitTime time.Duration) ConfigOption {
	return func(c *Config) {
		c.RetryCount = count
		c.RetryWaitTime = waitTime
	}
}
