package client

import (
	"context"
	"net/http"
)

// ListRoutingPolicies returns every routing policy in the org (base records,
// without rules). The list is unpaginated; callers filter client-side.
func (c *ScalingClient) ListRoutingPolicies(ctx context.Context) ([]RoutingPolicy, error) {
	data, err := doRequest[[]RoutingPolicy](c, ctx, http.MethodGet, "/v1/routing/policies", nil)
	if err != nil {
		return nil, err
	}
	return *data, nil
}

func (c *ScalingClient) GetRoutingPolicy(ctx context.Context, policyID string) (*RoutingPolicyWithRules, error) {
	return doRequest[RoutingPolicyWithRules](c, ctx, http.MethodGet, "/v1/routing/policies/"+policyID, nil)
}

func (c *ScalingClient) CreateRoutingPolicy(ctx context.Context, req CreateRoutingPolicyRequest) (*RoutingPolicyWithRules, error) {
	return doRequest[RoutingPolicyWithRules](c, ctx, http.MethodPost, "/v1/routing/policies", req)
}

func (c *ScalingClient) UpdateRoutingPolicy(ctx context.Context, policyID string, req UpdateRoutingPolicyRequest) (*RoutingPolicyWithRules, error) {
	return doRequest[RoutingPolicyWithRules](c, ctx, http.MethodPut, "/v1/routing/policies/"+policyID, req)
}

func (c *ScalingClient) DeleteRoutingPolicy(ctx context.Context, policyID string) error {
	_, err := doRequest[any](c, ctx, http.MethodDelete, "/v1/routing/policies/"+policyID, nil)
	return err
}
