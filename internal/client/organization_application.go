package client

import (
	"context"
	"fmt"
)

// OrgApplication represents an OAuth application under an organization.
type OrgApplication struct {
	ClientID       string `json:"client_id"`
	ClientSecret   string `json:"client_secret,omitempty"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	ApplicationURI string `json:"application_uri,omitempty"`
	RedirectURI    string `json:"redirect_uri,omitempty"`
	AvatarEmail    string `json:"avatar_email,omitempty"`
}

type CreateOrgApplicationRequest struct {
	Name           string `json:"name"`
	RedirectURI    string `json:"redirect_uri,omitempty"`
	ApplicationURI string `json:"application_uri,omitempty"`
	Description    string `json:"description,omitempty"`
	AvatarEmail    string `json:"avatar_email,omitempty"`
}

type UpdateOrgApplicationRequest struct {
	Name           string `json:"name"`
	RedirectURI    string `json:"redirect_uri"`
	ApplicationURI string `json:"application_uri"`
	Description    string `json:"description,omitempty"`
	AvatarEmail    string `json:"avatar_email,omitempty"`
}

func (c *Client) CreateOrgApplication(ctx context.Context, orgName string, req CreateOrgApplicationRequest) (*OrgApplication, error) {
	var app OrgApplication
	if err := c.post(ctx, fmt.Sprintf("organization/%s/applications", orgName), req, &app); err != nil {
		return nil, fmt.Errorf("creating application %q in org %q: %w", req.Name, orgName, err)
	}
	return &app, nil
}

func (c *Client) GetOrgApplication(ctx context.Context, orgName, clientID string) (*OrgApplication, error) {
	var app OrgApplication
	if err := c.get(ctx, fmt.Sprintf("organization/%s/applications/%s", orgName, clientID), &app); err != nil {
		return nil, fmt.Errorf("getting application %q in org %q: %w", clientID, orgName, err)
	}
	return &app, nil
}

func (c *Client) UpdateOrgApplication(ctx context.Context, orgName, clientID string, req UpdateOrgApplicationRequest) (*OrgApplication, error) {
	var app OrgApplication
	if err := c.put(ctx, fmt.Sprintf("organization/%s/applications/%s", orgName, clientID), req, &app); err != nil {
		return nil, fmt.Errorf("updating application %q in org %q: %w", clientID, orgName, err)
	}
	return &app, nil
}

func (c *Client) DeleteOrgApplication(ctx context.Context, orgName, clientID string) error {
	if err := c.delete(ctx, fmt.Sprintf("organization/%s/applications/%s", orgName, clientID)); err != nil {
		return fmt.Errorf("deleting application %q in org %q: %w", clientID, orgName, err)
	}
	return nil
}
