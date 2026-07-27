package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type SellCargoRequest struct {
	ShipSymbol string `json:"-"`
	Symbol     string `json:"symbol"`
	Units      int    `json:"units"`
}

type SellCargoResponse struct {
	Agent       ss.Agent       `json:"agent"`
	Transaction ss.Transaction `json:"transaction"`
	Cargo       ss.ShipCargo   `json:"cargo"`
}

// SellCargo sells units of a good from the ship's cargo at the current
// marketplace. The ship must be DOCKED at a waypoint with the Marketplace trait.
//
// POST /my/ships/{shipSymbol}/sell.
func (c *Client) SellCargo(
	ctx context.Context,
	req SellCargoRequest,
) (SellCargoResponse, error) {
	var res SellCargoResponse

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + req.ShipSymbol + "/sell",
		Body:      req,
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return SellCargoResponse{}, err
	}

	return res, nil
}
