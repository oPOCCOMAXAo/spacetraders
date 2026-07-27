package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type GetMyShipRequest struct {
	ShipSymbol string `json:"-"`
}

type GetMyShipResponse struct {
	Ship ss.Ship
}

// GetMyShip retrieves the full details of a single ship.
//
// GET /my/ships/{shipSymbol}.
func (c *Client) GetMyShip(
	ctx context.Context,
	req GetMyShipRequest,
) (GetMyShipResponse, error) {
	var ship ss.Ship

	err := c.Request(ctx, Request{
		Method:    http.MethodGet,
		Path:      "/my/ships/" + req.ShipSymbol,
		ResultRef: &ship,
	}, RequestOptions{})
	if err != nil {
		return GetMyShipResponse{}, err
	}

	return GetMyShipResponse{Ship: ship}, nil
}
