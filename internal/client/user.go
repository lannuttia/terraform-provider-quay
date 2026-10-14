package client

import (
	"context"
	"fmt"
)

// User represents a Quay user.
type User struct {
	Username  string `json:"username"`
	Email     string `json:"email,omitempty"`
	Verified  bool   `json:"verified,omitempty"`
	Anonymous bool   `json:"anonymous,omitempty"`
}

func (c *Client) GetCurrentUser(ctx context.Context) (*User, error) {
	var user User
	if err := c.get(ctx, "user/", &user); err != nil {
		return nil, fmt.Errorf("getting current user: %w", err)
	}
	return &user, nil
}

func (c *Client) GetUser(ctx context.Context, username string) (*User, error) {
	var user User
	if err := c.get(ctx, fmt.Sprintf("users/%s", username), &user); err != nil {
		return nil, fmt.Errorf("getting user %q: %w", username, err)
	}
	return &user, nil
}
