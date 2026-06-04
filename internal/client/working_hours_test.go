package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateWorkingHours(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}
		if r.URL.Path != "/v1/working-hours" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/working-hours")
		}

		var body CreateWorkingHoursRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.Name != "UK office hours" {
			t.Errorf("body.Name = %q, want %q", body.Name, "UK office hours")
		}
		if body.Timezone != "Europe/London" {
			t.Errorf("body.Timezone = %q, want %q", body.Timezone, "Europe/London")
		}
		if len(body.Windows) != 1 {
			t.Fatalf("len(body.Windows) = %d, want 1", len(body.Windows))
		}
		if len(body.Windows[0].Days) != 5 || body.Windows[0].Days[0] != 1 {
			t.Errorf("Windows[0].Days = %v, want [1 2 3 4 5]", body.Windows[0].Days)
		}
		if body.Windows[0].Start != "09:00" || body.Windows[0].End != "17:00" {
			t.Errorf("Windows[0] = %+v, want 09:00-17:00", body.Windows[0])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":       "wh_new",
				"orgId":    "org_1",
				"name":     "UK office hours",
				"timezone": "Europe/London",
				"windows": []map[string]any{
					{"days": []int{1, 2, 3, 4, 5}, "start": "09:00", "end": "17:00"},
				},
				"createdAt": "2026-06-04T10:00:00.000Z",
				"updatedAt": "2026-06-04T10:00:00.000Z",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.CreateWorkingHours(context.Background(), CreateWorkingHoursRequest{
		Name:     "UK office hours",
		Timezone: "Europe/London",
		Windows: []WorkingHoursWindow{
			{Days: []int{1, 2, 3, 4, 5}, Start: "09:00", End: "17:00"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "wh_new" {
		t.Errorf("ID = %q, want %q", result.ID, "wh_new")
	}
	if result.Timezone != "Europe/London" {
		t.Errorf("Timezone = %q, want %q", result.Timezone, "Europe/London")
	}
	if len(result.Windows) != 1 {
		t.Fatalf("len(Windows) = %d, want 1", len(result.Windows))
	}
}

func TestGetWorkingHours(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/working-hours/wh_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/working-hours/wh_123")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":       "wh_123",
				"orgId":    "org_1",
				"name":     "US office hours",
				"timezone": "America/New_York",
				"windows": []map[string]any{
					{"days": []int{1, 2, 3, 4, 5}, "start": "08:00", "end": "18:00"},
				},
				"createdAt": "2026-01-05T08:00:00.000Z",
				"updatedAt": "2026-03-20T16:45:00.000Z",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.GetWorkingHours(context.Background(), "wh_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "US office hours" {
		t.Errorf("Name = %q, want %q", result.Name, "US office hours")
	}
	if result.Timezone != "America/New_York" {
		t.Errorf("Timezone = %q, want %q", result.Timezone, "America/New_York")
	}
}

func TestUpdateWorkingHours(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("Method = %q, want PATCH", r.Method)
		}
		if r.URL.Path != "/v1/working-hours/wh_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/working-hours/wh_123")
		}

		var body UpdateWorkingHoursRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.Name != "Renamed" {
			t.Errorf("body.Name = %q, want %q", body.Name, "Renamed")
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"id":       "wh_123",
				"orgId":    "org_1",
				"name":     "Renamed",
				"timezone": "America/New_York",
				"windows": []map[string]any{
					{"days": []int{1, 2, 3, 4, 5}, "start": "08:00", "end": "18:00"},
				},
				"createdAt": "2026-01-05T08:00:00.000Z",
				"updatedAt": "2026-06-04T08:00:00.000Z",
			},
		})
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	result, err := c.UpdateWorkingHours(context.Background(), "wh_123", UpdateWorkingHoursRequest{
		Name:     "Renamed",
		Timezone: "America/New_York",
		Windows: []WorkingHoursWindow{
			{Days: []int{1, 2, 3, 4, 5}, Start: "08:00", End: "18:00"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "Renamed" {
		t.Errorf("Name = %q, want %q", result.Name, "Renamed")
	}
}

func TestDeleteWorkingHours(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Method = %q, want DELETE", r.Method)
		}
		if r.URL.Path != "/v1/working-hours/wh_123" {
			t.Errorf("Path = %q, want %q", r.URL.Path, "/v1/working-hours/wh_123")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "test-key")
	if err := c.DeleteWorkingHours(context.Background(), "wh_123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
