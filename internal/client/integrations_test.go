package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetInboundIntegration(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/integrations/inbound/int_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/integrations/inbound/int_123")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":              "int_123",
				"orgId":           "org_1",
				"name":            "Datadog",
				"componentId":     "cmp_1",
				"routingPolicyId": nil,
				"selectors": []map[string]any{
					{
						"matchers": []map[string]any{
							{"key": "service", "value": "payments"},
							{"key": "env", "value": "prod"},
						},
						"routingPolicyId": "rp_critical",
					},
					{
						"matchers":        []map[string]any{{"key": "env", "value": "staging"}},
						"routingPolicyId": "rp_low",
					},
				},
				"createdAt": "2026-01-05T08:00:00.000Z",
				"updatedAt": "2026-03-20T16:45:00.000Z",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.GetInboundIntegration(context.Background(), "int_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Datadog" || result.ComponentID != "cmp_1" {
		t.Errorf("got %q/%q, want Datadog/cmp_1", result.Name, result.ComponentID)
	}
	if len(result.Selectors) != 2 {
		t.Fatalf("len(Selectors) = %d, want 2", len(result.Selectors))
	}
	// Ordering is meaningful: the first row is the critical/prod one.
	if result.Selectors[0].RoutingPolicyID != "rp_critical" {
		t.Errorf("Selectors[0].RoutingPolicyID = %q, want rp_critical", result.Selectors[0].RoutingPolicyID)
	}
	if len(result.Selectors[0].Matchers) != 2 || result.Selectors[0].Matchers[0].Key != "service" {
		t.Errorf("Selectors[0].Matchers = %+v, unexpected", result.Selectors[0].Matchers)
	}
	if result.Selectors[1].RoutingPolicyID != "rp_low" {
		t.Errorf("Selectors[1].RoutingPolicyID = %q, want rp_low", result.Selectors[1].RoutingPolicyID)
	}
}

func TestSetInboundSelectors(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/v1/integrations/inbound/int_123/selectors" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/integrations/inbound/int_123/selectors")
		}

		var body SetSelectorsRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.Selectors == nil {
			t.Errorf("body.Selectors = nil, want a non-nil slice for full replacement")
		}
		if len(body.Selectors) != 1 {
			t.Fatalf("len(body.Selectors) = %d, want 1", len(body.Selectors))
		}
		if body.Selectors[0].RoutingPolicyID != "rp_new" {
			t.Errorf("Selectors[0].RoutingPolicyID = %q, want rp_new", body.Selectors[0].RoutingPolicyID)
		}
		if len(body.Selectors[0].Matchers) != 1 || body.Selectors[0].Matchers[0].Key != "team" {
			t.Errorf("Selectors[0].Matchers = %+v, unexpected", body.Selectors[0].Matchers)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":              "int_123",
				"orgId":           "org_1",
				"name":            "Datadog",
				"componentId":     "cmp_1",
				"routingPolicyId": nil,
				"selectors": []map[string]any{
					{
						"matchers":        []map[string]any{{"key": "team", "value": "payments"}},
						"routingPolicyId": "rp_new",
					},
				},
				"createdAt": "2026-01-05T08:00:00.000Z",
				"updatedAt": "2026-04-08T08:00:00.000Z",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.SetInboundSelectors(context.Background(), "int_123", SetSelectorsRequest{
		Selectors: []RoutingSelector{
			{
				Matchers:        []SelectorMatcher{{Key: "team", Value: "payments"}},
				RoutingPolicyID: "rp_new",
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Selectors) != 1 || result.Selectors[0].RoutingPolicyID != "rp_new" {
		t.Errorf("result.Selectors = %+v, unexpected", result.Selectors)
	}
}

func TestSetInboundSelectorsEmptyClears(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body SetSelectorsRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		// Clearing selectors must serialize as an empty array, not null.
		if body.Selectors == nil {
			t.Errorf("body.Selectors = nil, want an empty (non-nil) slice")
		}
		if len(body.Selectors) != 0 {
			t.Errorf("len(body.Selectors) = %d, want 0", len(body.Selectors))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":              "int_123",
				"orgId":           "org_1",
				"name":            "Datadog",
				"componentId":     "cmp_1",
				"routingPolicyId": nil,
				"selectors":       []map[string]any{},
				"createdAt":       "2026-01-05T08:00:00.000Z",
				"updatedAt":       "2026-04-08T08:00:00.000Z",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.SetInboundSelectors(context.Background(), "int_123", SetSelectorsRequest{
		Selectors: []RoutingSelector{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Selectors) != 0 {
		t.Errorf("result.Selectors = %+v, want empty", result.Selectors)
	}
}
