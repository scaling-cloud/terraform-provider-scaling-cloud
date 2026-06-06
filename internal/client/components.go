package client

import (
	"context"
	"net/http"
)

func (c *ScalingClient) GetComponent(ctx context.Context, componentID string) (*Component, error) {
	return doRequest[Component](c, ctx, http.MethodGet, "/v1/components/"+componentID, nil)
}

func (c *ScalingClient) CreateComponent(ctx context.Context, req CreateComponentRequest) (*Component, error) {
	return doRequest[Component](c, ctx, http.MethodPost, "/v1/components", req)
}

func (c *ScalingClient) UpdateComponent(ctx context.Context, componentID string, req UpdateComponentRequest) (*Component, error) {
	return doRequest[Component](c, ctx, http.MethodPatch, "/v1/components/"+componentID, req)
}

func (c *ScalingClient) DeleteComponent(ctx context.Context, componentID string) error {
	_, err := doRequest[any](c, ctx, http.MethodDelete, "/v1/components/"+componentID, nil)
	return err
}
