package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

type OrbitShipRequest struct {
	ShipSymbol string `json:"-"`
}

type OrbitShipResponse struct {
	Nav ss.ShipNav `json:"nav"`
}

// OrbitShip commands a ship into orbit (required to extract resources).
//
// POST /my/ships/{shipSymbol}/orbit.
func (c *Client) OrbitShip(
	ctx context.Context,
	req OrbitShipRequest,
) (OrbitShipResponse, error) {
	var res OrbitShipResponse

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + req.ShipSymbol + "/orbit",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return OrbitShipResponse{}, err
	}

	return res, nil
}
