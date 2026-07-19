package spacetraders

import (
	"context"
	"log/slog"

	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
)

var _ lifecycle.Servable = (*Client)(nil)

type Client struct {
	token string
	host  string

	logger *slog.Logger
}

func NewClient(
	config Config,
	logger *slog.Logger,
) (*Client, error) {
	return &Client{
		token:  config.Token,
		host:   config.Host,
		logger: logger,
	}, nil
}

func (c *Client) Serve(ctx context.Context) error {
	return nil
}
