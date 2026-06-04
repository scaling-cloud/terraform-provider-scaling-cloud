package client

import (
	"context"
	"net/http"
)

func (c *ScalingClient) GetWorkingHours(ctx context.Context, id string) (*WorkingHours, error) {
	return doRequest[WorkingHours](c, ctx, http.MethodGet, "/v1/working-hours/"+id, nil)
}

func (c *ScalingClient) CreateWorkingHours(ctx context.Context, req CreateWorkingHoursRequest) (*WorkingHours, error) {
	return doRequest[WorkingHours](c, ctx, http.MethodPost, "/v1/working-hours", req)
}

func (c *ScalingClient) UpdateWorkingHours(ctx context.Context, id string, req UpdateWorkingHoursRequest) (*WorkingHours, error) {
	return doRequest[WorkingHours](c, ctx, http.MethodPatch, "/v1/working-hours/"+id, req)
}

func (c *ScalingClient) DeleteWorkingHours(ctx context.Context, id string) error {
	_, err := doRequest[any](c, ctx, http.MethodDelete, "/v1/working-hours/"+id, nil)
	return err
}
