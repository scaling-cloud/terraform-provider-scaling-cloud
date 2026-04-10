package client

import (
	"context"
	"net/http"
)

func (c *ScalingClient) CreateOncallLayer(ctx context.Context, scheduleID string, req CreateLayerRequest) (*OncallLayer, error) {
	return doRequest[OncallLayer](c, ctx, http.MethodPost, "/v1/oncall/schedules/"+scheduleID+"/layers", req)
}

func (c *ScalingClient) UpdateOncallLayer(ctx context.Context, scheduleID, layerID string, req UpdateLayerRequest) (*OncallLayer, error) {
	return doRequest[OncallLayer](c, ctx, http.MethodPut, "/v1/oncall/schedules/"+scheduleID+"/layers/"+layerID, req)
}

func (c *ScalingClient) DeleteOncallLayer(ctx context.Context, scheduleID, layerID string) error {
	_, err := doRequest[any](c, ctx, http.MethodDelete, "/v1/oncall/schedules/"+scheduleID+"/layers/"+layerID, nil)
	return err
}
