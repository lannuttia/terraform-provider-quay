package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const defaultBaseURL = "https://quay.io"

// Client is a Quay API client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New creates a new Quay API client.
func New(baseURL, token string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	baseURL = strings.TrimRight(baseURL, "/")
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{},
	}
}

// BaseURL returns the base URL this client is configured to talk to.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// APIError represents an error response from the Quay API.
type APIError struct {
	StatusCode int
	Message    string
	ErrorType  string `json:"error_type"`
}

func (e *APIError) Error() string {
	if e.ErrorType != "" {
		return fmt.Sprintf("quay API error %d (%s): %s", e.StatusCode, e.ErrorType, e.Message)
	}
	return fmt.Sprintf("quay API error %d: %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error chain contains a 404 Not Found APIError.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

func (c *Client) newRequest(ctx context.Context, method, path string, body any) (*http.Request, error) {
	u, err := url.JoinPath(c.baseURL, "/api/v1/", path)
	if err != nil {
		return nil, fmt.Errorf("building URL: %w", err)
	}

	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshalling request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		apiErr := &APIError{StatusCode: resp.StatusCode}
		// Try to parse the error body
		var errBody struct {
			Message   string `json:"message"`
			ErrorType string `json:"error_type"`
		}
		if json.Unmarshal(respBody, &errBody) == nil {
			apiErr.Message = errBody.Message
			apiErr.ErrorType = errBody.ErrorType
		}
		if apiErr.Message == "" {
			apiErr.Message = string(respBody)
		}
		return apiErr
	}

	if out != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

func (c *Client) get(ctx context.Context, path string, out any) error {
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) post(ctx context.Context, path string, body, out any) error {
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) put(ctx context.Context, path string, body, out any) error {
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return err
	}
	return c.do(req, out)
}

func (c *Client) delete(ctx context.Context, path string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
