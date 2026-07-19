package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type SellCargoRequest struct {
	Symbol string `json:"symbol"`
	Units  int    `json:"units"`
}

// SellCargo sells units of a good from the ship's cargo at the current
// marketplace. The ship must be DOCKED at a waypoint with the Marketplace trait.
//
// POST /my/ships/{shipSymbol}/sell.
func (c *Client) SellCargo(
	ctx context.Context,
	shipSymbol string,
	req SellCargoRequest,
) (ss.Agent, ss.Transaction, ss.ShipCargo, error) {
	var res struct {
		Agent       ss.Agent       `json:"agent"`
		Transaction ss.Transaction `json:"transaction"`
		Cargo       ss.ShipCargo   `json:"cargo"`
	}

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + shipSymbol + "/sell",
		Body:      req,
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.Agent{}, ss.Transaction{}, ss.ShipCargo{}, err
	}

	return res.Agent, res.Transaction, res.Cargo, nil
}
