package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

// DockShip commands a ship to dock (required for market/contract actions).
//
// POST /my/ships/{shipSymbol}/dock.
func (c *Client) DockShip(
	ctx context.Context,
	shipSymbol string,
) (ss.ShipNav, error) {
	var res struct {
		Nav ss.ShipNav `json:"nav"`
	}

	err := c.Request(ctx, Request{
		Method:    http.MethodPost,
		Path:      "/my/ships/" + shipSymbol + "/dock",
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.ShipNav{}, err
	}

	return res.Nav, nil
}
