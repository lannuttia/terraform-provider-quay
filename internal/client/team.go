package client

import (
	"context"
	"fmt"
)

// Team represents a Quay organization team.
type Team struct {
	Name        string       `json:"name"`
	Role        string       `json:"role"`
	Description string       `json:"description,omitempty"`
	Members     []TeamMember `json:"members,omitempty"`
}

// TeamMember represents a member of a team.
type TeamMember struct {
	Name      string `json:"name"`
	Kind      string `json:"kind"` // "user" or "robot"
	IsRobot   bool   `json:"is_robot"`
	Invited   bool   `json:"invited,omitempty"`
	AvatarData map[string]any `json:"avatar,omitempty"`
}

type teamMembersResponse struct {
	Members []TeamMember `json:"members"`
	CanEdit bool         `json:"can_edit"`
}

type UpsertTeamRequest struct {
	Role        string `json:"role"`
	Description string `json:"description,omitempty"`
}

func (c *Client) UpsertOrgTeam(ctx context.Context, orgName, teamName string, req UpsertTeamRequest) (*Team, error) {
	var team Team
	if err := c.put(ctx, fmt.Sprintf("organization/%s/team/%s", orgName, teamName), req, &team); err != nil {
		return nil, fmt.Errorf("upserting team %q in org %q: %w", teamName, orgName, err)
	}
	team.Name = teamName
	return &team, nil
}

func (c *Client) GetOrgTeamMembers(ctx context.Context, orgName, teamName string, includePending bool) ([]TeamMember, error) {
	path := fmt.Sprintf("organization/%s/team/%s/members", orgName, teamName)
	if includePending {
		path += "?includePending=true"
	}
	var resp teamMembersResponse
	if err := c.get(ctx, path, &resp); err != nil {
		return nil, fmt.Errorf("getting members of team %q in org %q: %w", teamName, orgName, err)
	}
	return resp.Members, nil
}

func (c *Client) AddOrgTeamMember(ctx context.Context, orgName, teamName, memberName string) error {
	if err := c.put(ctx, fmt.Sprintf("organization/%s/team/%s/members/%s", orgName, teamName, memberName), nil, nil); err != nil {
		return fmt.Errorf("adding member %q to team %q in org %q: %w", memberName, teamName, orgName, err)
	}
	return nil
}

func (c *Client) RemoveOrgTeamMember(ctx context.Context, orgName, teamName, memberName string) error {
	if err := c.delete(ctx, fmt.Sprintf("organization/%s/team/%s/members/%s", orgName, teamName, memberName)); err != nil {
		return fmt.Errorf("removing member %q from team %q in org %q: %w", memberName, teamName, orgName, err)
	}
	return nil
}

func (c *Client) DeleteOrgTeam(ctx context.Context, orgName, teamName string) error {
	if err := c.delete(ctx, fmt.Sprintf("organization/%s/team/%s", orgName, teamName)); err != nil {
		return fmt.Errorf("deleting team %q in org %q: %w", teamName, orgName, err)
	}
	return nil
}
