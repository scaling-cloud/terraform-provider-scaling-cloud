package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewScalingClient(t *testing.T) {
	t.Parallel()

	c, err := NewScalingClient("https://api.scaling.cloud", "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.baseURL != "https://api.scaling.cloud" {
		t.Errorf("BaseURL = %q, want %q", c.baseURL, "https://api.scaling.cloud")
	}
}

func TestNewScalingClient_TrailingSlash(t *testing.T) {
	t.Parallel()

	c, err := NewScalingClient("https://api.scaling.cloud/", "test-key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.baseURL != "https://api.scaling.cloud" {
		t.Errorf("BaseURL = %q, want %q", c.baseURL, "https://api.scaling.cloud")
	}
}

func TestNewScalingClient_InvalidURL(t *testing.T) {
	t.Parallel()

	_, err := NewScalingClient("://invalid", "test-key")
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestNewScalingClient_RejectsHTTP(t *testing.T) {
	t.Parallel()

	_, err := NewScalingClient("http://api.scaling.cloud", "test-key")
	if err == nil {
		t.Fatal("expected error for non-HTTPS URL")
	}
}

func TestNewScalingClient_AllowsLocalhostHTTP(t *testing.T) {
	t.Parallel()

	_, err := NewScalingClient("http://localhost:8080", "test-key")
	if err != nil {
		t.Fatalf("unexpected error for localhost HTTP: %v", err)
	}
}

func TestDoRequest_Success(t *testing.T) {
	t.Parallel()

	type schedule struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("Authorization = %q, want %q", got, "Bearer test-key")
		}
		if got := r.URL.Path; got != "/v1/oncall/schedules/sched_123" {
			t.Errorf("Path = %q, want %q", got, "/v1/oncall/schedules/sched_123")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":   "sched_123",
				"name": "Primary",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := doRequest[schedule](c, context.Background(), http.MethodGet, "/v1/oncall/schedules/sched_123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "sched_123" {
		t.Errorf("ID = %q, want %q", result.ID, "sched_123")
	}
	if result.Name != "Primary" {
		t.Errorf("Name = %q, want %q", result.Name, "Primary")
	}
}

func TestDoRequest_404(t *testing.T) {
	t.Parallel()

	type empty struct{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"statusCode": 404,
			"type":       "invalid_request_error",
			"code":       "not_found",
			"message":    "Schedule not found",
			"requestId":  "req-abc",
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	_, err := doRequest[empty](c, context.Background(), http.MethodGet, "/v1/oncall/schedules/missing", nil)
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound = false, want true; error: %v", err)
	}
}

func TestDoRequest_401(t *testing.T) {
	t.Parallel()

	type empty struct{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"statusCode": 401,
			"type":       "not_authorized_error",
			"code":       "invalid_api_key",
			"message":    "Invalid API key",
			"requestId":  "req-def",
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "bad-key")
	_, err := doRequest[empty](c, context.Background(), http.MethodGet, "/v1/oncall/schedules", nil)
	if err == nil {
		t.Fatal("expected error for 401 response")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
	if apiErr.Code != "invalid_api_key" {
		t.Errorf("Code = %q, want %q", apiErr.Code, "invalid_api_key")
	}
}

func TestDoRequest_204NoContent(t *testing.T) {
	t.Parallel()

	type empty struct{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := doRequest[empty](c, context.Background(), http.MethodDelete, "/v1/oncall/schedules/sched_123", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for 204, got %v", result)
	}
}

func TestDoRequest_WithBody(t *testing.T) {
	t.Parallel()

	type createReq struct {
		Name     string `json:"name"`
		Timezone string `json:"timezone"`
	}
	type schedule struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want %q", ct, "application/json")
		}

		var body createReq
		json.NewDecoder(r.Body).Decode(&body)
		if body.Name != "Primary" {
			t.Errorf("body.Name = %q, want %q", body.Name, "Primary")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":   "sched_new",
				"name": "Primary",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := doRequest[schedule](c, context.Background(), http.MethodPost, "/v1/oncall/schedules", createReq{
		Name:     "Primary",
		Timezone: "America/New_York",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "sched_new" {
		t.Errorf("ID = %q, want %q", result.ID, "sched_new")
	}
}
