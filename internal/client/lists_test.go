package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// listServer replies to any GET with a {"data": <body>} envelope so the bare
// list endpoints (which return an array under "data") can be exercised.
func listServer(t *testing.T, wantPath, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %q, want GET", r.Method)
		}
		if r.URL.Path != wantPath {
			t.Errorf("Path = %q, want %q", r.URL.Path, wantPath)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":` + body + `}`))
	}))
}

func TestListEscalationPolicies(t *testing.T) {
	t.Parallel()
	server := listServer(t, "/v1/escalation/policies", `[
		{"id":"ep_1","orgId":"org_1","name":"Critical","description":null,"createdAt":"t","updatedAt":"t"},
		{"id":"ep_2","orgId":"org_1","name":"Low","description":"low pri","createdAt":"t","updatedAt":"t"}
	]`)
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "k")
	got, err := c.ListEscalationPolicies(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].Name != "Critical" || got[1].ID != "ep_2" {
		t.Errorf("got %+v, unexpected", got)
	}
}

func TestListRoutingPolicies(t *testing.T) {
	t.Parallel()
	server := listServer(t, "/v1/routing/policies", `[
		{"id":"rp_1","orgId":"org_1","name":"Default","description":null,"isDefault":true,"createdAt":"t","updatedAt":"t"}
	]`)
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "k")
	got, err := c.ListRoutingPolicies(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Default" || !got[0].IsDefault {
		t.Errorf("got %+v, unexpected", got)
	}
}

func TestListOncallSchedules(t *testing.T) {
	t.Parallel()
	server := listServer(t, "/v1/oncall/schedules", `[
		{"id":"sch_1","orgId":"org_1","name":"Primary","description":null,"timezone":"Europe/London","createdAt":"t","updatedAt":"t"}
	]`)
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "k")
	got, err := c.ListOncallSchedules(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Primary" || got[0].Timezone != "Europe/London" {
		t.Errorf("got %+v, unexpected", got)
	}
}

func TestListWorkingHours(t *testing.T) {
	t.Parallel()
	server := listServer(t, "/v1/working-hours", `[
		{"id":"wh_1","orgId":"org_1","name":"Business Hours","timezone":"Europe/London","windows":[],"createdAt":"t","updatedAt":"t"}
	]`)
	defer server.Close()

	c, _ := NewScalingClient(server.URL, "k")
	got, err := c.ListWorkingHours(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "Business Hours" || got[0].Timezone != "Europe/London" {
		t.Errorf("got %+v, unexpected", got)
	}
}
