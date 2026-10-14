package client

import (
	"context"
	"fmt"
)

// PrototypeDelegate is a user or team that receives a permission prototype.
type PrototypeDelegate struct {
	Kind string `json:"kind"` // "user" or "team"
	Name string `json:"name"`
}

// PrototypeActivatingUser is the user whose repo creation activates the prototype.
type PrototypeActivatingUser struct {
	Name string `json:"name"`
}

// Prototype represents a default permission rule for an organization.
type Prototype struct {
	ID              string                   `json:"id"`
	Role            string                   `json:"role"`
	Delegate        PrototypeDelegate        `json:"delegate"`
	ActivatingUser  *PrototypeActivatingUser `json:"activating_user,omitempty"`
}

type prototypesResponse struct {
	Prototypes []Prototype `json:"prototypes"`
}

type CreatePrototypeRequest struct {
	Role           string                   `json:"role"`
	Delegate       PrototypeDelegate        `json:"delegate"`
	ActivatingUser *PrototypeActivatingUser `json:"activating_user,omitempty"`
}

type UpdatePrototypeRequest struct {
	Role string `json:"role"`
}

func (c *Client) ListOrgPrototypes(ctx context.Context, orgName string) ([]Prototype, error) {
	var resp prototypesResponse
	if err := c.get(ctx, fmt.Sprintf("organization/%s/prototypes", orgName), &resp); err != nil {
		return nil, fmt.Errorf("listing prototypes for org %q: %w", orgName, err)
	}
	return resp.Prototypes, nil
}

func (c *Client) CreateOrgPrototype(ctx context.Context, orgName string, req CreatePrototypeRequest) (*Prototype, error) {
	var proto Prototype
	if err := c.post(ctx, fmt.Sprintf("organization/%s/prototypes", orgName), req, &proto); err != nil {
		return nil, fmt.Errorf("creating prototype in org %q: %w", orgName, err)
	}
	return &proto, nil
}

func (c *Client) GetOrgPrototype(ctx context.Context, orgName, protoID string) (*Prototype, error) {
	protos, err := c.ListOrgPrototypes(ctx, orgName)
	if err != nil {
		return nil, err
	}
	for _, p := range protos {
		if p.ID == protoID {
			return &p, nil
		}
	}
	return nil, &APIError{StatusCode: 404, Message: fmt.Sprintf("prototype %q not found in org %q", protoID, orgName)}
}

func (c *Client) UpdateOrgPrototype(ctx context.Context, orgName, protoID string, req UpdatePrototypeRequest) (*Prototype, error) {
	var proto Prototype
	if err := c.put(ctx, fmt.Sprintf("organization/%s/prototypes/%s", orgName, protoID), req, &proto); err != nil {
		return nil, fmt.Errorf("updating prototype %q in org %q: %w", protoID, orgName, err)
	}
	return &proto, nil
}

func (c *Client) DeleteOrgPrototype(ctx context.Context, orgName, protoID string) error {
	if err := c.delete(ctx, fmt.Sprintf("organization/%s/prototypes/%s", orgName, protoID)); err != nil {
		return fmt.Errorf("deleting prototype %q in org %q: %w", protoID, orgName, err)
	}
	return nil
}
