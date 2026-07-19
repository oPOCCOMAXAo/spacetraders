package spacetraders

import (
	"context"
	"net/http"

	ss "github.com/opoccomaxao/spacetraders/pkg/services/spacetraders/structs"
)

// GetMyShip retrieves the full details of a single ship.
//
// GET /my/ships/{shipSymbol}.
func (c *Client) GetMyShip(
	ctx context.Context,
	shipSymbol string,
) (ss.Ship, error) {
	var res ss.Ship

	err := c.Request(ctx, Request{
		Method:    http.MethodGet,
		Path:      "/my/ships/" + shipSymbol,
		ResultRef: &res,
	}, RequestOptions{})
	if err != nil {
		return ss.Ship{}, err
	}

	return res, nil
}
