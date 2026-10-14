package client

import (
	"context"
	"fmt"
)

// Repository represents a Quay repository.
type Repository struct {
	Namespace   string `json:"namespace"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"is_public"`
	Kind        string `json:"kind,omitempty"`
}

type CreateRepositoryRequest struct {
	Repository  string `json:"repository"`
	Namespace   string `json:"namespace,omitempty"`
	Visibility  string `json:"visibility"`
	Description string `json:"description"`
	RepoKind    string `json:"repo_kind,omitempty"`
}

type UpdateRepositoryRequest struct {
	Description string `json:"description"`
}

type ChangeVisibilityRequest struct {
	Visibility string `json:"visibility"`
}

func (c *Client) CreateRepository(ctx context.Context, req CreateRepositoryRequest) (*Repository, error) {
	var repo Repository
	if err := c.post(ctx, "repository", req, &repo); err != nil {
		return nil, fmt.Errorf("creating repository %q: %w", req.Repository, err)
	}
	return c.GetRepository(ctx, fmt.Sprintf("%s/%s", req.Namespace, req.Repository))
}

func (c *Client) GetRepository(ctx context.Context, repository string) (*Repository, error) {
	var repo Repository
	if err := c.get(ctx, fmt.Sprintf("repository/%s", repository), &repo); err != nil {
		return nil, fmt.Errorf("getting repository %q: %w", repository, err)
	}
	return &repo, nil
}

func (c *Client) UpdateRepository(ctx context.Context, repository string, req UpdateRepositoryRequest) (*Repository, error) {
	if err := c.put(ctx, fmt.Sprintf("repository/%s", repository), req, nil); err != nil {
		return nil, fmt.Errorf("updating repository %q: %w", repository, err)
	}
	return c.GetRepository(ctx, repository)
}

func (c *Client) ChangeRepositoryVisibility(ctx context.Context, repository string, req ChangeVisibilityRequest) error {
	if err := c.post(ctx, fmt.Sprintf("repository/%s/changevisibility", repository), req, nil); err != nil {
		return fmt.Errorf("changing visibility of repository %q: %w", repository, err)
	}
	return nil
}

func (c *Client) DeleteRepository(ctx context.Context, repository string) error {
	if err := c.delete(ctx, fmt.Sprintf("repository/%s", repository)); err != nil {
		return fmt.Errorf("deleting repository %q: %w", repository, err)
	}
	return nil
}
