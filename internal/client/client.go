package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const maxResponseBodySize = 10 * 1024 * 1024 // 10 MB

type ScalingClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewScalingClient(baseURL, apiKey string) (*ScalingClient, error) {
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base_url: %w", err)
	}

	if parsed.Scheme != "https" && parsed.Hostname() != "localhost" && parsed.Hostname() != "127.0.0.1" {
		return nil, fmt.Errorf("base_url must use https scheme, got %q (http is only allowed for localhost)", parsed.Scheme)
	}

	transport := &http.Transport{
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	return &ScalingClient{
		baseURL: strings.TrimRight(parsed.String(), "/"),
		httpClient: &http.Client{
			Transport: &authTransport{
				authHeader: "Bearer " + apiKey,
				userAgent:  "terraform-provider-scaling-cloud/0.1.0",
				base: &retryTransport{
					base:       transport,
					maxRetries: 3,
					baseDelay:  1 * time.Second,
				},
			},
			Timeout: 30 * time.Second,
		},
	}, nil
}

func doRequest[T any](c *ScalingClient, ctx context.Context, method, path string, body any) (*T, error) {
	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	reqURL := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBodySize))
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			return nil, fmt.Errorf("HTTP %d: unexpected error response", resp.StatusCode)
		}
		return nil, &apiErr
	}

	var envelope struct {
		Data T `json:"data"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &envelope.Data, nil
}
