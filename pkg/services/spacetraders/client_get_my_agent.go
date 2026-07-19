package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

// GetMyAgent fetches the caller's agent details.
//
// GET /my/agent.
func (c *Client) GetMyAgent(ctx context.Context) (ss.Agent, error) {
	var res ss.Agent

	err := c.Request(ctx, Request{
		Method:    http.MethodGet,
		Path:      "/my/agent",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.Agent{}, err
	}

	return res, nil
}
