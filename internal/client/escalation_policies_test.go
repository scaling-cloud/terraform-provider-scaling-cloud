package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateEscalationPolicy_withConditions(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/v1/escalation/policies" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/escalation/policies")
		}

		var body CreateEscalationPolicyRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if len(body.Steps) != 2 {
			t.Fatalf("len(body.Steps) = %d, want 2", len(body.Steps))
		}
		if body.Steps[0].Condition == nil {
			t.Fatalf("Steps[0].Condition = nil, want a condition")
		}
		if body.Steps[0].Condition.WorkingHoursID != "wh_uk" || body.Steps[0].Condition.When != "during" {
			t.Errorf("Steps[0].Condition = %+v, want wh_uk/during", *body.Steps[0].Condition)
		}
		if body.Steps[1].Condition != nil {
			t.Errorf("Steps[1].Condition = %+v, want nil", body.Steps[1].Condition)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":        "ep_new",
				"orgId":     "org_1",
				"name":      "Follow the sun",
				"createdAt": "2026-06-04T10:00:00.000Z",
				"updatedAt": "2026-06-04T10:00:00.000Z",
				"steps": []map[string]any{
					{
						"id":                   "step_1",
						"position":             1,
						"targetType":           "schedule",
						"targetId":             "sch_uk",
						"escalateAfterSeconds": 300,
						"condition":            map[string]any{"workingHoursId": "wh_uk", "when": "during"},
						"createdAt":            "2026-06-04T10:00:00.000Z",
						"updatedAt":            "2026-06-04T10:00:00.000Z",
					},
					{
						"id":                   "step_2",
						"position":             2,
						"targetType":           "schedule",
						"targetId":             "sch_us",
						"escalateAfterSeconds": 600,
						"condition":            nil,
						"createdAt":            "2026-06-04T10:00:00.000Z",
						"updatedAt":            "2026-06-04T10:00:00.000Z",
					},
				},
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.CreateEscalationPolicy(context.Background(), CreateEscalationPolicyRequest{
		Name: "Follow the sun",
		Steps: []EscalationStepInput{
			{
				Position:             1,
				TargetType:           "schedule",
				TargetID:             "sch_uk",
				EscalateAfterSeconds: 300,
				Condition:            &WorkingHoursCondition{WorkingHoursID: "wh_uk", When: "during"},
			},
			{
				Position:             2,
				TargetType:           "schedule",
				TargetID:             "sch_us",
				EscalateAfterSeconds: 600,
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Steps) != 2 {
		t.Fatalf("len(Steps) = %d, want 2", len(result.Steps))
	}
	if result.Steps[0].Condition == nil {
		t.Fatalf("Steps[0].Condition = nil, want a condition")
	}
	if result.Steps[0].Condition.WorkingHoursID != "wh_uk" || result.Steps[0].Condition.When != "during" {
		t.Errorf("Steps[0].Condition = %+v, want wh_uk/during", *result.Steps[0].Condition)
	}
	if result.Steps[1].Condition != nil {
		t.Errorf("Steps[1].Condition = %+v, want nil", result.Steps[1].Condition)
	}
}

func TestGetEscalationPolicy_withConditions(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/escalation/policies/ep_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/escalation/policies/ep_123")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":        "ep_123",
				"orgId":     "org_1",
				"name":      "Follow the sun",
				"createdAt": "2026-01-05T08:00:00.000Z",
				"updatedAt": "2026-03-20T16:45:00.000Z",
				"steps": []map[string]any{
					{
						"id":                   "step_1",
						"position":             1,
						"targetType":           "schedule",
						"targetId":             "sch_us",
						"escalateAfterSeconds": 300,
						"condition":            map[string]any{"workingHoursId": "wh_us", "when": "outside"},
						"createdAt":            "2026-01-05T08:00:00.000Z",
						"updatedAt":            "2026-01-05T08:00:00.000Z",
					},
				},
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.GetEscalationPolicy(context.Background(), "ep_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Steps) != 1 {
		t.Fatalf("len(Steps) = %d, want 1", len(result.Steps))
	}
	if result.Steps[0].Condition == nil {
		t.Fatalf("Steps[0].Condition = nil, want a condition")
	}
	if result.Steps[0].Condition.When != "outside" {
		t.Errorf("Steps[0].Condition.When = %q, want outside", result.Steps[0].Condition.When)
	}
}
