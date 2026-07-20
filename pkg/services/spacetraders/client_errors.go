package spacetraders

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	pkgerr "github.com/opoccomaxao/gopkg/pkg/errors"
	pkgerrors "github.com/pkg/errors"
)

// APIError mirrors the SpaceTraders error envelope:
//
//	{"error": {"message": "...", "code": 4xxx, "data": {...}}}
type APIError struct {
	Message string          `json:"message"`
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
}

func (e *APIError) Error() string {
	if e == nil {
		return ""
	}

	return "spacetraders api error " + strconv.Itoa(e.Code) + ": " + e.Message
}

// errorEnvelope is the outer wrapper around APIError.
type errorEnvelope struct {
	Error *APIError `json:"error,omitempty"`
}

func (*Client) noLogErrors(int, []byte) {}

func (c *Client) verboseLogErrors(statusCode int, body []byte) {
	attrs := []any{
		slog.Int("status_code", statusCode),
	}

	var env errorEnvelope

	err := json.Unmarshal(body, &env)
	if err == nil && env.Error != nil {
		attrs = append(attrs,
			slog.Int("code", env.Error.Code),
			slog.String("message", env.Error.Message),
			slog.String("data", string(env.Error.Data)),
		)
	} else {
		attrs = append(attrs, slog.String("body", string(body)))
	}

	c.logger.Error("spacetraders request error", attrs...)
}

// mapError translates an HTTP status + parsed APIError into a wrapped sentinel
// from gopkg/pkg/errors so callers can use errors.Is.
func (*Client) mapError(statusCode int, apiErr *APIError) error {
	msg := "unexpected status code"
	if apiErr != nil && apiErr.Message != "" {
		msg = apiErr.Message
	}

	var sentinel error

	switch {
	case statusCode == http.StatusUnauthorized:
		sentinel = pkgerr.ErrInvalidAuth
	case statusCode == http.StatusNotFound:
		sentinel = pkgerr.ErrNotFound
	case statusCode == http.StatusTooManyRequests:
		sentinel = pkgerr.ErrLimitExceeded
	case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
		sentinel = pkgerr.ErrInvalidParam
	case statusCode >= http.StatusInternalServerError:
		sentinel = pkgerr.ErrRequestFailed
	default:
		sentinel = pkgerr.ErrFailed
	}

	return pkgerrors.Wrap(sentinel, msg)
}
