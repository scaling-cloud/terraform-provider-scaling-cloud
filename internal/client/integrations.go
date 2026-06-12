package client

import (
	"context"
	"net/http"
)

// ListInboundIntegrations returns every inbound integration in the org. The
// list is unpaginated; callers filter client-side (for example, by name) to
// resolve a single integration.
func (c *ScalingClient) ListInboundIntegrations(ctx context.Context) ([]InboundIntegration, error) {
	data, err := doRequest[struct {
		Integrations []InboundIntegration `json:"integrations"`
	}](c, ctx, http.MethodGet, "/v1/integrations/inbound", nil)
	if err != nil {
		return nil, err
	}
	return data.Integrations, nil
}

func (c *ScalingClient) GetInboundIntegration(ctx context.Context, integrationID string) (*InboundIntegration, error) {
	return doRequest[InboundIntegration](c, ctx, http.MethodGet, "/v1/integrations/inbound/"+integrationID, nil)
}

// SetInboundSelectors replaces the integration's ordered Routing Selectors
// wholesale (ADR-0039). An empty (non-nil) selectors slice clears them.
func (c *ScalingClient) SetInboundSelectors(ctx context.Context, integrationID string, req SetSelectorsRequest) (*InboundIntegration, error) {
	if req.Selectors == nil {
		req.Selectors = []RoutingSelector{}
	}
	return doRequest[InboundIntegration](c, ctx, http.MethodPut, "/v1/integrations/inbound/"+integrationID+"/selectors", req)
}
