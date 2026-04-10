package client

import (
	"context"
	"net/http"
)

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
