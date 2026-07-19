package spacetraders

import "time"

const (
	defaultHost     = "https://api.spacetraders.io/v2"
	rateLimitPerSec = 2.0
	rateLimitBurst  = 1

	httpTimeout = 30 * time.Second

	maxRetries       = 5
	retryBaseBackoff = 500 * time.Millisecond
	retryMaxBackoff  = 8 * time.Second

	defaultPageSize = 20

	authScheme = "Bearer "
)
