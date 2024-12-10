package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Version fetches response from /versions endpoint.
func (c *Client) Version(projectKey string) ([]*Version, error) {
	path := fmt.Sprintf("/project/%s/versions", projectKey)
	res, err := c.GetV2(context.Background(), path, nil)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, ErrEmptyResponse
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return nil, formatUnexpectedResponse(res)
	}

	var out []*Version

	err = json.NewDecoder(res.Body).Decode(&out)

	return out, err
}
