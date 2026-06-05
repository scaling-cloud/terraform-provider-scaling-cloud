package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateComponent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/v1/components" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/components")
		}

		var body CreateComponentRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.Name != "Payments API" {
			t.Errorf("body.Name = %q, want %q", body.Name, "Payments API")
		}
		if len(body.Aliases) != 2 || body.Aliases[0] != "payments" || body.Aliases[1] != "checkout" {
			t.Errorf("body.Aliases = %v, want [payments checkout]", body.Aliases)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":                "cmp_new",
				"orgId":             "org_1",
				"name":              "Payments API",
				"description":       nil,
				"aliases":           []string{"payments", "checkout"},
				"operationalStatus": "operational",
				"maintenanceWindow": nil,
				"createdAt":         "2026-04-07T10:00:00.000Z",
				"updatedAt":         "2026-04-07T10:00:00.000Z",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.CreateComponent(context.Background(), CreateComponentRequest{
		Name:    "Payments API",
		Aliases: []string{"payments", "checkout"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "cmp_new" {
		t.Errorf("ID = %q, want %q", result.ID, "cmp_new")
	}
	if len(result.Aliases) != 2 || result.Aliases[0] != "payments" {
		t.Errorf("Aliases = %v, want [payments checkout]", result.Aliases)
	}
}

func TestGetComponent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/components/cmp_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/components/cmp_123")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":                "cmp_123",
				"orgId":             "org_1",
				"name":              "Payments API",
				"description":       "Handles payments",
				"aliases":           []string{"payments"},
				"operationalStatus": "operational",
				"maintenanceWindow": nil,
				"createdAt":         "2026-01-05T08:00:00.000Z",
				"updatedAt":         "2026-03-20T16:45:00.000Z",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.GetComponent(context.Background(), "cmp_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Description == nil || *result.Description != "Handles payments" {
		t.Errorf("Description = %v, want Handles payments", result.Description)
	}
	if len(result.Aliases) != 1 || result.Aliases[0] != "payments" {
		t.Errorf("Aliases = %v, want [payments]", result.Aliases)
	}
}

func TestUpdateComponent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Method = %q, want PATCH", r.Method)
		}
		if r.URL.Path != "/v1/components/cmp_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/components/cmp_123")
		}

		var body UpdateComponentRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		// Full-replacement: aliases is always sent (non-nil), even when empty.
		if body.Aliases == nil {
			t.Errorf("body.Aliases = nil, want a non-nil slice for full replacement")
		}
		if len(body.Aliases) != 0 {
			t.Errorf("body.Aliases = %v, want empty slice", body.Aliases)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":                "cmp_123",
				"orgId":             "org_1",
				"name":              "Renamed",
				"description":       nil,
				"aliases":           []string{},
				"operationalStatus": "operational",
				"maintenanceWindow": nil,
				"createdAt":         "2026-01-05T08:00:00.000Z",
				"updatedAt":         "2026-04-08T08:00:00.000Z",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.UpdateComponent(context.Background(), "cmp_123", UpdateComponentRequest{
		Name:    "Renamed",
		Aliases: []string{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Renamed" {
		t.Errorf("Name = %q, want %q", result.Name, "Renamed")
	}
	if len(result.Aliases) != 0 {
		t.Errorf("Aliases = %v, want empty", result.Aliases)
	}
}

func TestDeleteComponent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/v1/components/cmp_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/components/cmp_123")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	if err := c.DeleteComponent(context.Background(), "cmp_123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
