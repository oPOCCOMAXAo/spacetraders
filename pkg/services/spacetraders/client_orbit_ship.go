package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

// OrbitShip commands a ship into orbit (required to extract resources).
//
// POST /my/ships/{shipSymbol}/orbit.
func (c *Client) OrbitShip(
	ctx context.Context,
	shipSymbol string,
) (ss.ShipNav, error) {
	var res struct {
		Nav ss.ShipNav `json:"nav"`
	}

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + shipSymbol + "/orbit",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.ShipNav{}, err
	}

	return res.Nav, nil
}
