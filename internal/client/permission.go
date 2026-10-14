package client

import (
	"context"
	"fmt"
)

// Permission represents a repository permission for a user or team.
type Permission struct {
	Role string `json:"role"`
	Name string `json:"name,omitempty"`
}

type SetPermissionRequest struct {
	Role string `json:"role"`
}

// Team permissions

func (c *Client) GetRepoTeamPermission(ctx context.Context, repository, teamName string) (*Permission, error) {
	var perm Permission
	if err := c.get(ctx, fmt.Sprintf("repository/%s/permissions/team/%s", repository, teamName), &perm); err != nil {
		return nil, fmt.Errorf("getting team permission for %q on %q: %w", teamName, repository, err)
	}
	return &perm, nil
}

func (c *Client) SetRepoTeamPermission(ctx context.Context, repository, teamName string, req SetPermissionRequest) (*Permission, error) {
	var perm Permission
	if err := c.put(ctx, fmt.Sprintf("repository/%s/permissions/team/%s", repository, teamName), req, &perm); err != nil {
		return nil, fmt.Errorf("setting team permission for %q on %q: %w", teamName, repository, err)
	}
	return &perm, nil
}

func (c *Client) DeleteRepoTeamPermission(ctx context.Context, repository, teamName string) error {
	if err := c.delete(ctx, fmt.Sprintf("repository/%s/permissions/team/%s", repository, teamName)); err != nil {
		return fmt.Errorf("deleting team permission for %q on %q: %w", teamName, repository, err)
	}
	return nil
}

// User permissions

func (c *Client) GetRepoUserPermission(ctx context.Context, repository, username string) (*Permission, error) {
	var perm Permission
	if err := c.get(ctx, fmt.Sprintf("repository/%s/permissions/user/%s", repository, username), &perm); err != nil {
		return nil, fmt.Errorf("getting user permission for %q on %q: %w", username, repository, err)
	}
	return &perm, nil
}

func (c *Client) SetRepoUserPermission(ctx context.Context, repository, username string, req SetPermissionRequest) (*Permission, error) {
	var perm Permission
	if err := c.put(ctx, fmt.Sprintf("repository/%s/permissions/user/%s", repository, username), req, &perm); err != nil {
		return nil, fmt.Errorf("setting user permission for %q on %q: %w", username, repository, err)
	}
	return &perm, nil
}

func (c *Client) DeleteRepoUserPermission(ctx context.Context, repository, username string) error {
	if err := c.delete(ctx, fmt.Sprintf("repository/%s/permissions/user/%s", repository, username)); err != nil {
		return fmt.Errorf("deleting user permission for %q on %q: %w", username, repository, err)
	}
	return nil
}
