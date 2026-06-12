package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListUsers(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Method = %q, want GET", r.Method)
		}
		if r.URL.Path != "/v1/users" {
			t.Errorf("Path = %q, want /v1/users", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"data":{"users":[
			{"id":"user_1","firstName":"Jane","lastName":"Smith","email":"jane@example.com"},
			{"id":"user_2","firstName":null,"lastName":null,"email":null}
		]}}`))
	}))
	defer srv.Close()

	c, _ := NewScalingClient(srv.URL, "k")
	got, err := c.ListUsers(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0].ID != "user_1" {
		t.Fatalf("got %+v, unexpected", got)
	}
	if got[0].Email == nil || *got[0].Email != "jane@example.com" {
		t.Errorf("got[0].Email = %v, want jane@example.com", got[0].Email)
	}
	if got[1].Email != nil {
		t.Errorf("got[1].Email = %v, want nil", got[1].Email)
	}
}
