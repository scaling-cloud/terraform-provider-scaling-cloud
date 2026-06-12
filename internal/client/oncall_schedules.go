package client

import (
	"context"
	"net/http"
)

// ListOncallSchedules returns every on-call schedule in the org (base records,
// without layers). The list is unpaginated; callers filter client-side.
func (c *ScalingClient) ListOncallSchedules(ctx context.Context) ([]OncallSchedule, error) {
	data, err := doRequest[[]OncallSchedule](c, ctx, http.MethodGet, "/v1/oncall/schedules", nil)
	if err != nil {
		return nil, err
	}
	return *data, nil
}

func (c *ScalingClient) GetOncallSchedule(ctx context.Context, scheduleID string) (*OncallScheduleWithLayers, error) {
	return doRequest[OncallScheduleWithLayers](c, ctx, http.MethodGet, "/v1/oncall/schedules/"+scheduleID, nil)
}

func (c *ScalingClient) CreateOncallSchedule(ctx context.Context, req CreateScheduleRequest) (*OncallSchedule, error) {
	return doRequest[OncallSchedule](c, ctx, http.MethodPost, "/v1/oncall/schedules", req)
}

func (c *ScalingClient) UpdateOncallSchedule(ctx context.Context, scheduleID string, req UpdateScheduleRequest) (*OncallSchedule, error) {
	return doRequest[OncallSchedule](c, ctx, http.MethodPut, "/v1/oncall/schedules/"+scheduleID, req)
}

func (c *ScalingClient) DeleteOncallSchedule(ctx context.Context, scheduleID string) error {
	_, err := doRequest[any](c, ctx, http.MethodDelete, "/v1/oncall/schedules/"+scheduleID, nil)
	return err
}
