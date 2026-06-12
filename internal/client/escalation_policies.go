package client

import (
	"context"
	"net/http"
)

// ListEscalationPolicies returns every escalation policy in the org (base
// records, without steps). The list is unpaginated; callers filter client-side.
func (c *ScalingClient) ListEscalationPolicies(ctx context.Context) ([]EscalationPolicy, error) {
	data, err := doRequest[[]EscalationPolicy](c, ctx, http.MethodGet, "/v1/escalation/policies", nil)
	if err != nil {
		return nil, err
	}
	return *data, nil
}

func (c *ScalingClient) GetEscalationPolicy(ctx context.Context, policyID string) (*EscalationPolicyWithSteps, error) {
	return doRequest[EscalationPolicyWithSteps](c, ctx, http.MethodGet, "/v1/escalation/policies/"+policyID, nil)
}

func (c *ScalingClient) CreateEscalationPolicy(ctx context.Context, req CreateEscalationPolicyRequest) (*EscalationPolicyWithSteps, error) {
	return doRequest[EscalationPolicyWithSteps](c, ctx, http.MethodPost, "/v1/escalation/policies", req)
}

func (c *ScalingClient) UpdateEscalationPolicy(ctx context.Context, policyID string, req UpdateEscalationPolicyRequest) (*EscalationPolicyWithSteps, error) {
	return doRequest[EscalationPolicyWithSteps](c, ctx, http.MethodPut, "/v1/escalation/policies/"+policyID, req)
}

func (c *ScalingClient) DeleteEscalationPolicy(ctx context.Context, policyID string) error {
	_, err := doRequest[any](c, ctx, http.MethodDelete, "/v1/escalation/policies/"+policyID, nil)
	return err
}
