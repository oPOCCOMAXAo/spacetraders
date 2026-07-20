package spacetraders

import (
	"context"
	"log/slog"
	"net/http"

	pkgerr "github.com/opoccomaxao/gopkg/pkg/errors"
	"github.com/opoccomaxao/gopkg/pkg/services/lifecycle"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/time/rate"
)

var _ lifecycle.Servable = (*Client)(nil)

// Client is a REST client for the SpaceTraders v2 API.
type Client struct {
	token string
	host  string

	logger    *slog.Logger
	client    *http.Client
	rateLimit *rate.Limiter

	logErrors func(statusCode int, body []byte)
}

func NewClient(
	config Config,
	logger *slog.Logger,
) (*Client, error) {
	if config.Token == "" {
		return nil, pkgerrors.WithStack(pkgerr.ErrInvalidParam)
	}

	host := config.Host
	if host == "" {
		host = defaultHost
	}

	res := &Client{
		token: config.Token,
		host:  host,

		logger: logger,
		client: &http.Client{
			Timeout: httpTimeout,
		},
	}
	res.rateLimit = res.newRateLimiter()

	if config.VerboseErrors {
		res.logErrors = res.verboseLogErrors
	} else {
		res.logErrors = res.noLogErrors
	}

	return res, nil
}

func (*Client) Serve(context.Context) error {
	return nil
}
