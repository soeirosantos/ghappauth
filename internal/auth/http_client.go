package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/avast/retry-go/v4"
	"ghappauth/internal/types"
)

// HTTPClient wraps the standard http.Client with retry logic and common functionality
type HTTPClient struct {
	client *http.Client
	config *HTTPClientConfig
}

// HTTPClientConfig holds configuration for the HTTP client
type HTTPClientConfig struct {
	Timeout         time.Duration
	MaxRetries      uint
	RetryDelay      time.Duration
	BackoffMultiplier float64
	UserAgent       string
}

// DefaultHTTPClientConfig returns default configuration for the HTTP client
func DefaultHTTPClientConfig() *HTTPClientConfig {
	return &HTTPClientConfig{
		Timeout:         30 * time.Second,
		MaxRetries:      3,
		RetryDelay:      1 * time.Second,
		BackoffMultiplier: 2.0,
		UserAgent:       "ghappauth/1.0",
	}
}

// NewHTTPClient creates a new HTTP client with the given configuration
func NewHTTPClient(config *HTTPClientConfig) *HTTPClient {
	if config == nil {
		config = DefaultHTTPClientConfig()
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
	}
}

// RequestConfig holds configuration for individual requests
type RequestConfig struct {
	Method      string
	URL         string
	AuthToken   string
	Accept      string
	Body        io.Reader
	ExpectedStatus int
}

// RetryableError represents an error that should trigger a retry
type RetryableError struct {
	StatusCode int
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable status code: %d", e.StatusCode)
}

// doRequest performs an HTTP request with retry logic and common error handling
func (c *HTTPClient) doRequest(ctx context.Context, config *RequestConfig) (*http.Response, error) {
	var resp *http.Response
	
	err := retry.Do(
		func() error {
			req, err := http.NewRequestWithContext(ctx, config.Method, config.URL, config.Body)
			if err != nil {
				return fmt.Errorf("failed to create request: %w", err)
			}

			if config.AuthToken != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.AuthToken))
			}
			if config.Accept != "" {
				req.Header.Set("Accept", config.Accept)
			} else {
				req.Header.Set("Accept", "application/vnd.github.v3+json")
			}
			req.Header.Set("User-Agent", c.config.UserAgent)

			response, err := c.client.Do(req)
			if err != nil {
				return fmt.Errorf("failed to make request: %w", err)
			}

			if shouldRetry(response.StatusCode) {
				response.Body.Close()
				return &RetryableError{StatusCode: response.StatusCode}
			}

			resp = response
			return nil
		},
		retry.Attempts(c.config.MaxRetries),
		retry.Delay(c.config.RetryDelay),
		retry.DelayType(retry.BackOffDelay),
		retry.LastErrorOnly(true),
		retry.RetryIf(func(err error) bool {
			_, ok := err.(*RetryableError)
			return ok
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("request failed after %d attempts: %w", c.config.MaxRetries, err)
	}

	return resp, nil
}

// DoRequest performs an HTTP request and decodes the JSON response
func (c *HTTPClient) DoRequest(ctx context.Context, config *RequestConfig, result interface{}) error {
	resp, err := c.doRequest(ctx, config)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if config.ExpectedStatus != 0 && resp.StatusCode != config.ExpectedStatus {
		var apiError types.GitHubAPIError
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return fmt.Errorf("GitHub API error: status %d (failed to decode error response: %w)", resp.StatusCode, err)
		}
		return fmt.Errorf("GitHub API error: %s (status: %d)", apiError.Message, resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// shouldRetry determines if a status code should trigger a retry
func shouldRetry(statusCode int) bool {
	return statusCode >= 500 || statusCode == 429 || statusCode == 408
}



 