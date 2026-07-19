package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

// FulfillContract fulfills a contract whose delivery terms are all met,
// paying the on-fulfillment credits.
//
// POST /my/contracts/{contractId}/fulfill.
func (c *Client) FulfillContract(
	ctx context.Context,
	contractID string,
) (ss.Agent, ss.Contract, error) {
	var res struct {
		Agent    ss.Agent    `json:"agent"`
		Contract ss.Contract `json:"contract"`
	}

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/contracts/" + contractID + "/fulfill",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.Agent{}, ss.Contract{}, err
	}

	return res.Agent, res.Contract, nil
}
