package client

import "time"

// NewTestClient creates a client with zero retry delays, for use in tests only.
func NewTestClient(baseURL, token string) *Client {
	return New(baseURL, token).withRetryDelay(func(int) time.Duration { return 0 })
}
