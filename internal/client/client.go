// Package client provides an HTTP client for the eDME API with retry logic,
// token management, and keepalive support.
package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"terraform-provider-dme/internal/api"
)

const (
	// Default timeout for each API call.
	defaultTimeout = 30 * time.Second
	// Number of retries on transient failures.
	maxRetries = 2
	// API base path prefix.
	apiBasePath = "/rest"
)

// Config holds the configuration for the eDME API client.
type Config struct {
	Endpoint string // eDME management IP and port, e.g. "https://10.0.0.1:26335"
	UserName string // eDME northbound user name
	Password string // eDME northbound user password
}

// Client is the HTTP client for communicating with the eDME API.
type Client struct {
	endpoint    string
	userName    string
	password    string
	httpClient  *http.Client
	mu          sync.Mutex // Protects token and tokenExpiry
	token       string     // Current auth token (accessSession)
	tokenExpiry time.Time  // When the current token expires
}

// NewClient creates a new eDME API client.
func NewClient(cfg Config) *Client {
	return &Client{
		endpoint: cfg.Endpoint,
		userName: cfg.UserName,
		password: cfg.Password,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

// Authenticate calls the session API to obtain a new auth token.
// PUT /rest/plat/smapp/v1/sessions
func (c *Client) Authenticate(ctx context.Context) error {
	req := api.AuthRequest{
		GrantType: "password",
		UserName:  c.userName,
		Value:     c.password,
	}

	var resp api.AuthResponse
	if err := c.doRequest(ctx, http.MethodPut, "/plat/smapp/v1/sessions", req, &resp, false); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	c.mu.Lock()
	c.token = resp.AccessSession
	// Set expiry with a 60-second buffer to avoid edge-case race conditions.
	c.tokenExpiry = time.Now().Add(time.Duration(resp.Expires-60) * time.Second)
	c.mu.Unlock()

	return nil
}

// GetToken returns the current auth token, refreshing it if expired.
func (c *Client) GetToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	if c.token == "" || time.Now().After(c.tokenExpiry) {
		c.mu.Unlock()
		if err := c.Authenticate(ctx); err != nil {
			return "", err
		}
		c.mu.Lock()
	}
	token := c.token
	c.mu.Unlock()
	return token, nil
}

// doRequest executes an HTTP request with the given method, path, and body.
// If withToken is true, the request includes the X-Auth-Token header.
// It implements 2 retries on failure with a 30-second per-call timeout.
func (c *Client) doRequest(ctx context.Context, method, path string, body, result interface{}, withToken bool) error {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Second * time.Duration(attempt))
		}

		err := c.singleRequest(ctx, method, path, body, result, withToken)
		if err == nil {
			return nil
		}
		lastErr = err
	}
	return lastErr
}

// singleRequest performs a single HTTP request attempt.
func (c *Client) singleRequest(ctx context.Context, method, path string, body, result interface{}, withToken bool) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	url := c.endpoint + apiBasePath + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Charset", "utf8")
	req.Header.Set("Content-Type", "application/json")

	if withToken {
		token, err := c.GetToken(ctx)
		if err != nil {
			return fmt.Errorf("get token: %w", err)
		}
		req.Header.Set("X-Auth-Token", token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	// Read the body for error reporting.
	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("api error %d %s %s: %s", resp.StatusCode, method, path, string(respBody))
	}

	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

// Do is the public method for making API calls with token authentication.
func (c *Client) Do(ctx context.Context, method, path string, body, result interface{}) error {
	return c.doRequest(ctx, method, path, body, result, true)
}
