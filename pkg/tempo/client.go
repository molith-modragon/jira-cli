package tempo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ankitpokhrel/jira-cli/pkg/jira"
)

const (
	// Tempo API base path
	tempoAPIPath = "/rest/tempo-teams/1/team"
)

// Client is a Tempo client that uses the Jira client for authentication.
type Client struct {
	jiraClient *jira.Client
}

// NewClient creates a new Tempo client using the provided Jira client.
func NewClient(jiraClient *jira.Client) *Client {
	return &Client{
		jiraClient: jiraClient,
	}
}

// GetTeams retrieves all teams from Tempo.
func (c *Client) GetTeams() ([]*Team, error) {
	ctx := context.Background()
	
	// Make request to Tempo API
	res, err := c.jiraClient.Get(ctx, tempoAPIPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get teams: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get teams: %s", res.Status)
	}

	var teams []*Team
	if err := json.NewDecoder(res.Body).Decode(&teams); err != nil {
		return nil, fmt.Errorf("failed to decode teams response: %w", err)
	}

	return teams, nil
}