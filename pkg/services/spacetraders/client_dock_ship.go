package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type DockShipRequest struct {
	ShipSymbol string `json:"-"`
}

type DockShipResponse struct {
	Nav ss.ShipNav `json:"nav"`
}

// DockShip commands a ship to dock (required for market/contract actions).
//
// POST /my/ships/{shipSymbol}/dock.
func (c *Client) DockShip(
	ctx context.Context,
	req DockShipRequest,
) (DockShipResponse, error) {
	var res DockShipResponse

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + req.ShipSymbol + "/dock",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return DockShipResponse{}, err
	}

	return res, nil
}
