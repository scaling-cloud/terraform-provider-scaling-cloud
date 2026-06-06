package client

import (
	"context"
	"net/http"
)

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
