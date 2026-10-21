package client

import "time"

// NewTestClient creates a client with zero retry delays and no HTTP timeout,
// for use in tests only.
func NewTestClient(baseURL, token string) *Client {
	return New(baseURL, token, 0).withRetryDelay(func(int) time.Duration { return 0 })
}
