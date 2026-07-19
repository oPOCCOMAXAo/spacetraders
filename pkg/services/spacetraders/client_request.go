package spacetraders

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	pkgerr "github.com/opoccomaxao/gopkg/pkg/errors"
	"github.com/opoccomaxao/gopkg/pkg/services/logger"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/time/rate"
)

const (
	backoffBase    = 2
	pageSizeDigits = 2
)

// Request describes a single REST call.
//
// Path is appended verbatim to the configured host (e.g. "/my/agent" or
// "/my/ships/SYM-1/navigate"). Path parameters must already be interpolated
// by the caller.
//
// Body, when non-nil, is JSON-encoded and the request becomes a POST/PUT.
// ResultRef, when non-nil, receives the decoded inner "data" object of the
// success envelope. MetaRef, when non-nil, receives the decoded "meta" object
// (used by paginated list endpoints).
type Request struct {
	Method    string
	Path      string
	Query     url.Values
	Body      any
	ResultRef any
	MetaRef   any
}

type RequestOptions struct {
	LogRequest  bool
	LogResponse bool
}

// Request executes a single REST call against the SpaceTraders API.
//
// It applies the proactive rate limiter, retries 429 (honouring Retry-After)
// and 5xx responses with exponential backoff up to maxRetries times, decodes
// the standard success/error envelopes, and maps API errors to
// gopkg/pkg/errors sentinels so callers can use errors.Is.
func (c *Client) Request(
	ctx context.Context,
	req Request,
	opts RequestOptions,
) error {
	err := ctx.Err()
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	err = c.rateLimit.Wait(ctx)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	fullURL := c.host + req.Path
	if len(req.Query) > 0 {
		fullURL += "?" + req.Query.Encode()
	}

	bodyReader, bodyBytes, err := c.encodeBody(req.Body)
	if err != nil {
		return err
	}

	log := logger.NewLogRecord(c.logger)
	defer log.Info("request")

	if !opts.LogRequest && !opts.LogResponse {
		log.Discard()
	}

	if opts.LogRequest {
		c.logRequest(log, req.Method, fullURL, bodyBytes)
	}

	statusCode, respBody, err := c.doWithRetry(
		ctx,
		req.Method,
		fullURL,
		bodyReader,
		req.Body != nil,
	)
	if err != nil {
		return err
	}

	if opts.LogResponse {
		log.AddAttrs(
			slog.Int("status_code", statusCode),
			slog.String("response", string(respBody)),
		)
	}

	if statusCode < http.StatusOK || statusCode > http.StatusNoContent {
		c.logErrors(statusCode, respBody)

		return pkgerrors.WithStack(mapError(statusCode, c.parseError(respBody)))
	}

	return c.decodeResult(respBody, req)
}

// encodeBody marshals the request body (if any) and returns a replayable
// reader plus the raw bytes.
func (*Client) encodeBody(body any) (io.Reader, []byte, error) {
	if body == nil {
		return http.NoBody, nil, nil
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, nil, pkgerrors.WithStack(err)
	}

	return bytes.NewReader(bodyBytes), bodyBytes, nil
}

func (*Client) logRequest(log *logger.LogRecord, method, fullURL string, bodyBytes []byte) {
	log.AddAttrs(
		slog.String("method", method),
		slog.String("url", fullURL),
	)

	if len(bodyBytes) > 0 {
		log.AddAttrs(slog.String("request_body", string(bodyBytes)))
	}
}

// decodeResult unmarshals the success envelope into ResultRef/MetaRef.
func (c *Client) decodeResult(respBody []byte, req Request) error {
	if req.ResultRef == nil && req.MetaRef == nil {
		return nil
	}

	// 204 No Content or empty body: nothing to decode.
	if len(respBody) == 0 {
		return nil
	}

	var env responseEnvelope

	err := json.Unmarshal(respBody, &env)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	err = c.decodeData(env.Data, req.ResultRef)
	if err != nil {
		return err
	}

	return c.decodeMeta(env.Meta, req.MetaRef)
}

// decodeData copies the envelope's raw "data" field into result when present.
func (c *Client) decodeData(raw json.RawMessage, result any) error {
	if result == nil || len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	err := json.Unmarshal(raw, result)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

// decodeMeta copies the envelope's raw "meta" field into the caller's
// MetaRef when both are present.
func (c *Client) decodeMeta(raw json.RawMessage, metaRef any) error {
	if metaRef == nil || len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	err := json.Unmarshal(raw, metaRef)
	if err != nil {
		return pkgerrors.WithStack(err)
	}

	return nil
}

// doWithRetry performs the HTTP call, retrying on 429 and 5xx with bounded
// exponential backoff. It returns the final HTTP status code, the fully-read
// response body, and any terminal error. The response body is always closed
// before returning.
func (c *Client) doWithRetry(
	ctx context.Context,
	method, fullURL string,
	body io.Reader,
	hasBody bool,
) (int, []byte, error) {
	// Snapshot the request body so retries can replay it.
	bodyBytes, err := c.snapshotBody(body, hasBody)
	if err != nil {
		return 0, nil, err
	}

	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		err := ctx.Err()
		if err != nil {
			if lastErr == nil {
				lastErr = pkgerrors.WithStack(err)
			}

			break
		}

		statusCode, respBody, retryAble, retryAfter, callErr := c.doOnce(
			ctx,
			method,
			fullURL,
			bodyBytes,
			hasBody,
		)
		switch {
		case callErr != nil:
			lastErr = callErr

			c.sleepBackoff(ctx, attempt)
		case retryAble:
			lastErr = pkgerrors.WithStack(mapError(statusCode, c.parseError(respBody)))
			c.sleepFor(ctx, retryAfter, attempt)
		default:
			return statusCode, respBody, nil
		}
	}

	if lastErr == nil {
		lastErr = pkgerrors.WithStack(pkgerr.ErrRequestFailed)
	}

	return 0, nil, lastErr
}

// snapshotBody reads the request body once so it can be replayed across
// retries.
func (*Client) snapshotBody(body io.Reader, hasBody bool) ([]byte, error) {
	if !hasBody {
		return nil, nil
	}

	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, pkgerrors.WithStack(err)
	}

	return bodyBytes, nil
}

// doOnce performs a single HTTP attempt. It returns the status code, the
// response body, whether the call should be retried (429/5xx), the parsed
// Retry-After duration (429 only), and any transport-level error.
func (c *Client) doOnce(
	ctx context.Context,
	method, fullURL string,
	bodyBytes []byte,
	hasBody bool,
) (int, []byte, bool, time.Duration, error) {
	var reqBody io.Reader = http.NoBody
	if hasBody {
		reqBody = bytes.NewReader(bodyBytes)
	}

	httpReq, reqErr := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if reqErr != nil {
		return 0, nil, false, 0, pkgerrors.WithStack(reqErr)
	}

	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("Authorization", authScheme+c.token)

	if hasBody {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, doErr := c.client.Do(httpReq)
	if doErr != nil {
		return 0, nil, true, 0, pkgerrors.WithStack(doErr)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return resp.StatusCode, nil, true, 0, pkgerrors.WithStack(readErr)
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		return resp.StatusCode, respBody, true, parseRetryAfter(resp), nil
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		return resp.StatusCode, respBody, true, 0, nil
	}

	return resp.StatusCode, respBody, false, 0, nil
}

func (c *Client) parseError(body []byte) *APIError {
	var env errorEnvelope

	// Best-effort: a malformed error body is not itself useful to report;
	// the caller still has the HTTP status to map.
	_ = json.Unmarshal(body, &env)

	return env.Error
}

// parseRetryAfter reads the Retry-After header (whole seconds per the
// SpaceTraders spec) and returns it as a Duration, or 0 if absent/invalid.
func parseRetryAfter(resp *http.Response) time.Duration {
	v := resp.Header.Get("Retry-After")
	if v == "" {
		return 0
	}

	secs, err := strconv.Atoi(v)
	if err != nil || secs <= 0 {
		return 0
	}

	return time.Duration(secs) * time.Second
}

// sleepFor sleeps for the given Retry-After duration when non-zero, otherwise
// falls back to exponential backoff for the current attempt.
func (c *Client) sleepFor(ctx context.Context, retryAfter time.Duration, attempt int) {
	if retryAfter > 0 {
		c.sleep(ctx, retryAfter)

		return
	}

	c.sleepBackoff(ctx, attempt)
}

// sleepBackoff sleeps for an exponentially-growing duration, capped at
// retryMaxBackoff.
func (c *Client) sleepBackoff(ctx context.Context, attempt int) {
	d := retryBaseBackoff
	for range attempt {
		d *= backoffBase
	}

	if d > retryMaxBackoff {
		d = retryMaxBackoff
	}

	c.sleep(ctx, d)
}

func (c *Client) sleep(ctx context.Context, d time.Duration) {
	if d <= 0 {
		return
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}

// newRateLimiter builds the default proactive limiter used by NewClient.
func newRateLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Limit(rateLimitPerSec), rateLimitBurst)
}
