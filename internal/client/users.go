package client

import (
	"context"
	"net/http"
)

// ListUsers returns every user in the org. The list is unpaginated; callers
// filter client-side (for example, by email) to resolve a single user.
func (c *ScalingClient) ListUsers(ctx context.Context) ([]User, error) {
	data, err := doRequest[struct {
		Users []User `json:"users"`
	}](c, ctx, http.MethodGet, "/v1/users", nil)
	if err != nil {
		return nil, err
	}
	return data.Users, nil
}
