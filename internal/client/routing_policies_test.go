package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateRoutingPolicy(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/v1/routing/policies" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/routing/policies")
		}

		var body CreateRoutingPolicyRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.Name != "Payments routing" {
			t.Errorf("body.Name = %q, want %q", body.Name, "Payments routing")
		}
		if len(body.Rules) != 4 {
			t.Fatalf("len(body.Rules) = %d, want 4", len(body.Rules))
		}
		if body.Rules[0].Severity != "critical" || body.Rules[0].Outcome != "incident" {
			t.Errorf("rule[0] = %+v, want critical/incident", body.Rules[0])
		}
		if body.Rules[0].EscalationPolicyID == nil || *body.Rules[0].EscalationPolicyID != "esc_1" {
			t.Errorf("rule[0].EscalationPolicyID = %v, want esc_1", body.Rules[0].EscalationPolicyID)
		}
		if body.Rules[1].EscalationPolicyID != nil {
			t.Errorf("rule[1].EscalationPolicyID = %v, want nil", body.Rules[1].EscalationPolicyID)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":        "rp_new",
				"orgId":     "org_1",
				"name":      "Payments routing",
				"isDefault": false,
				"createdAt": "2026-04-07T10:00:00.000Z",
				"updatedAt": "2026-04-07T10:00:00.000Z",
				"rules": []map[string]any{
					{"severity": "critical", "outcome": "incident", "escalationPolicyId": "esc_1"},
					{"severity": "high", "outcome": "provisional_page", "escalationPolicyId": nil},
					{"severity": "medium", "outcome": "notification", "escalationPolicyId": nil},
					{"severity": "low", "outcome": "drop", "escalationPolicyId": nil},
				},
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	escID := "esc_1"
	result, err := c.CreateRoutingPolicy(context.Background(), CreateRoutingPolicyRequest{
		Name: "Payments routing",
		Rules: []RoutingRuleInput{
			{Severity: "critical", Outcome: "incident", EscalationPolicyID: &escID},
			{Severity: "high", Outcome: "provisional_page"},
			{Severity: "medium", Outcome: "notification"},
			{Severity: "low", Outcome: "drop"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "rp_new" {
		t.Errorf("ID = %q, want %q", result.ID, "rp_new")
	}
	if result.IsDefault {
		t.Errorf("IsDefault = true, want false")
	}
	if len(result.Rules) != 4 {
		t.Fatalf("len(Rules) = %d, want 4", len(result.Rules))
	}
	if result.Rules[0].EscalationPolicyID == nil || *result.Rules[0].EscalationPolicyID != "esc_1" {
		t.Errorf("Rules[0].EscalationPolicyID = %v, want esc_1", result.Rules[0].EscalationPolicyID)
	}
	if result.Rules[1].EscalationPolicyID != nil {
		t.Errorf("Rules[1].EscalationPolicyID = %v, want nil", result.Rules[1].EscalationPolicyID)
	}
}

func TestGetRoutingPolicy(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/routing/policies/rp_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/routing/policies/rp_123")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":        "rp_123",
				"orgId":     "org_1",
				"name":      "Default routing",
				"isDefault": true,
				"createdAt": "2026-01-05T08:00:00.000Z",
				"updatedAt": "2026-03-20T16:45:00.000Z",
				"rules": []map[string]any{
					{"severity": "critical", "outcome": "incident", "escalationPolicyId": nil},
					{"severity": "high", "outcome": "incident", "escalationPolicyId": nil},
					{"severity": "medium", "outcome": "notification", "escalationPolicyId": nil},
					{"severity": "low", "outcome": "notification", "escalationPolicyId": nil},
				},
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.GetRoutingPolicy(context.Background(), "rp_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.IsDefault {
		t.Errorf("IsDefault = false, want true")
	}
	if len(result.Rules) != 4 {
		t.Fatalf("len(Rules) = %d, want 4", len(result.Rules))
	}
}

func TestUpdateRoutingPolicy(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Method = %q, want PUT", r.Method)
		}
		if r.URL.Path != "/v1/routing/policies/rp_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/routing/policies/rp_123")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":        "rp_123",
				"orgId":     "org_1",
				"name":      "Renamed",
				"isDefault": false,
				"createdAt": "2026-01-05T08:00:00.000Z",
				"updatedAt": "2026-04-08T08:00:00.000Z",
				"rules": []map[string]any{
					{"severity": "critical", "outcome": "incident", "escalationPolicyId": nil},
					{"severity": "high", "outcome": "incident", "escalationPolicyId": nil},
					{"severity": "medium", "outcome": "notification", "escalationPolicyId": nil},
					{"severity": "low", "outcome": "notification", "escalationPolicyId": nil},
				},
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.UpdateRoutingPolicy(context.Background(), "rp_123", UpdateRoutingPolicyRequest{
		Name: "Renamed",
		Rules: []RoutingRuleInput{
			{Severity: "critical", Outcome: "incident"},
			{Severity: "high", Outcome: "incident"},
			{Severity: "medium", Outcome: "notification"},
			{Severity: "low", Outcome: "notification"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Renamed" {
		t.Errorf("Name = %q, want %q", result.Name, "Renamed")
	}
}

func TestDeleteRoutingPolicy(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/v1/routing/policies/rp_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/routing/policies/rp_123")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	if err := c.DeleteRoutingPolicy(context.Background(), "rp_123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
