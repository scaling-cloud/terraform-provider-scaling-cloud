package client

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"net/http"
	"strconv"
	"time"
)

const maxRetryAfterSeconds = 60

type authTransport struct {
	authHeader string
	userAgent  string
	base       http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Set("Authorization", t.authHeader)
	r.Header.Set("User-Agent", t.userAgent)
	if r.Body != nil {
		r.Header.Set("Content-Type", "application/json")
	}
	return t.base.RoundTrip(r)
}

type retryTransport struct {
	base       http.RoundTripper
	maxRetries int
	baseDelay  time.Duration
}

func (t *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for attempt := range t.maxRetries + 1 {
		if attempt > 0 && req.GetBody != nil {
			req.Body, err = req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("resetting request body for retry: %w", err)
			}
		}

		resp, err = t.base.RoundTrip(req)

		if err != nil {
			if isRetryableNetworkError(err) && attempt < t.maxRetries {
				if sleepErr := sleepWithContext(req.Context(), backoff(t.baseDelay, attempt)); sleepErr != nil {
					return nil, sleepErr
				}
				continue
			}
			return nil, err
		}

		if !isRetryableStatus(resp.StatusCode) || attempt == t.maxRetries {
			return resp, nil
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			if ra := resp.Header.Get("Retry-After"); ra != "" {
				if seconds, parseErr := strconv.Atoi(ra); parseErr == nil {
					if seconds > maxRetryAfterSeconds {
						seconds = maxRetryAfterSeconds
					}
					resp.Body.Close()
					if sleepErr := sleepWithContext(req.Context(), time.Duration(seconds)*time.Second); sleepErr != nil {
						return nil, sleepErr
					}
					continue
				}
			}
		}

		resp.Body.Close()
		if sleepErr := sleepWithContext(req.Context(), backoff(t.baseDelay, attempt)); sleepErr != nil {
			return nil, sleepErr
		}
	}

	return resp, err
}

func isRetryableStatus(status int) bool {
	switch status {
	case http.StatusTooManyRequests,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		return true
	}
	return false
}

func isRetryableNetworkError(err error) bool {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	return false
}

func backoff(base time.Duration, attempt int) time.Duration {
	delay := base * time.Duration(1<<attempt)
	if delay > 30*time.Second {
		delay = 30 * time.Second
	}
	return time.Duration(rand.Int64N(int64(delay) + 1))
}

func sleepWithContext(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return nil
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
