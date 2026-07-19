package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

// AcceptContract accepts an offered contract, paying the upfront credits.
//
// POST /my/contracts/{contractId}/accept.
func (c *Client) AcceptContract(
	ctx context.Context,
	contractID string,
) (ss.Agent, ss.Contract, error) {
	var res struct {
		Agent    ss.Agent    `json:"agent"`
		Contract ss.Contract `json:"contract"`
	}

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/contracts/" + contractID + "/accept",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.Agent{}, ss.Contract{}, err
	}

	return res.Agent, res.Contract, nil
}
