package client

import (
	"context"
	"fmt"
)

// Robot represents a Quay robot account.
type Robot struct {
	Name                  string         `json:"name"`
	Token                 string         `json:"token,omitempty"`
	Description           string         `json:"description,omitempty"`
	UnstructuredMetadata  map[string]any `json:"unstructured_metadata,omitempty"`
}

type CreateRobotRequest struct {
	Description          string         `json:"description,omitempty"`
	UnstructuredMetadata map[string]any `json:"unstructured_metadata,omitempty"`
}

func (c *Client) CreateOrgRobot(ctx context.Context, orgName, shortName string, req CreateRobotRequest) (*Robot, error) {
	var robot Robot
	if err := c.put(ctx, fmt.Sprintf("organization/%s/robots/%s", orgName, shortName), req, &robot); err != nil {
		return nil, fmt.Errorf("creating robot %q in org %q: %w", shortName, orgName, err)
	}
	return &robot, nil
}

func (c *Client) GetOrgRobot(ctx context.Context, orgName, shortName string) (*Robot, error) {
	var robot Robot
	if err := c.get(ctx, fmt.Sprintf("organization/%s/robots/%s", orgName, shortName), &robot); err != nil {
		return nil, fmt.Errorf("getting robot %q in org %q: %w", shortName, orgName, err)
	}
	return &robot, nil
}

func (c *Client) DeleteOrgRobot(ctx context.Context, orgName, shortName string) error {
	if err := c.delete(ctx, fmt.Sprintf("organization/%s/robots/%s", orgName, shortName)); err != nil {
		return fmt.Errorf("deleting robot %q in org %q: %w", shortName, orgName, err)
	}
	return nil
}

func (c *Client) RegenerateOrgRobotToken(ctx context.Context, orgName, shortName string) (*Robot, error) {
	var robot Robot
	if err := c.post(ctx, fmt.Sprintf("organization/%s/robots/%s/regenerate", orgName, shortName), nil, &robot); err != nil {
		return nil, fmt.Errorf("regenerating token for robot %q in org %q: %w", shortName, orgName, err)
	}
	return &robot, nil
}

func (c *Client) CreateUserRobot(ctx context.Context, shortName string, req CreateRobotRequest) (*Robot, error) {
	var robot Robot
	if err := c.put(ctx, fmt.Sprintf("user/robots/%s", shortName), req, &robot); err != nil {
		return nil, fmt.Errorf("creating user robot %q: %w", shortName, err)
	}
	return &robot, nil
}

func (c *Client) GetUserRobot(ctx context.Context, shortName string) (*Robot, error) {
	var robot Robot
	if err := c.get(ctx, fmt.Sprintf("user/robots/%s", shortName), &robot); err != nil {
		return nil, fmt.Errorf("getting user robot %q: %w", shortName, err)
	}
	return &robot, nil
}

func (c *Client) DeleteUserRobot(ctx context.Context, shortName string) error {
	if err := c.delete(ctx, fmt.Sprintf("user/robots/%s", shortName)); err != nil {
		return fmt.Errorf("deleting user robot %q: %w", shortName, err)
	}
	return nil
}

func (c *Client) RegenerateUserRobotToken(ctx context.Context, shortName string) (*Robot, error) {
	var robot Robot
	if err := c.post(ctx, fmt.Sprintf("user/robots/%s/regenerate", shortName), nil, &robot); err != nil {
		return nil, fmt.Errorf("regenerating token for user robot %q: %w", shortName, err)
	}
	return &robot, nil
}
