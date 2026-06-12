package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// componentsPageLimit is the page size requested when walking the paginated
// components list. The endpoint caps limit at 100.
const componentsPageLimit = 100

// maxComponentPages bounds the pagination walk so a misbehaving server that
// never clears nextCursor cannot loop forever.
const maxComponentPages = 1000

// ListComponents returns every component in the org. Unlike the other list
// endpoints, components are cursor-paginated, so this walks every page until
// nextCursor is empty, then callers filter the full set client-side.
func (c *ScalingClient) ListComponents(ctx context.Context) ([]Component, error) {
	var all []Component
	cursor := ""

	for page := 0; page < maxComponentPages; page++ {
		path := fmt.Sprintf("/v1/components?limit=%d", componentsPageLimit)
		if cursor != "" {
			path += "&cursor=" + url.QueryEscape(cursor)
		}

		data, err := doRequest[struct {
			Components []Component `json:"components"`
			NextCursor *string     `json:"nextCursor"`
		}](c, ctx, http.MethodGet, path, nil)
		if err != nil {
			return nil, err
		}

		all = append(all, data.Components...)
		if data.NextCursor == nil || *data.NextCursor == "" {
			return all, nil
		}
		cursor = *data.NextCursor
	}

	return nil, fmt.Errorf("listing components: exceeded %d pages without exhausting the cursor", maxComponentPages)
}

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
