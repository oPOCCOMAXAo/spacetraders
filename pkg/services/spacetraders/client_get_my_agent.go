package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type GetMyAgentRequest struct{}

type GetMyAgentResponse struct {
	Agent ss.Agent
}

// GetMyAgent fetches the caller's agent details.
//
// GET /my/agent.
func (c *Client) GetMyAgent(
	ctx context.Context,
	_ GetMyAgentRequest,
) (GetMyAgentResponse, error) {
	var agent ss.Agent

	err := c.Request(ctx, Request{
		Method:    http.MethodGet,
		Path:      "/my/agent",
		ResultRef: &agent,
	}, RequestOptions{})
	if err != nil {
		return GetMyAgentResponse{}, err
	}

	return GetMyAgentResponse{Agent: agent}, nil
}
