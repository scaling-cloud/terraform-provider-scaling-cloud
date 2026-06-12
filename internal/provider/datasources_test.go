package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func registeredDataSourceTypeNames(t *testing.T) map[string]bool {
	t.Helper()
	p := New("test")()
	names := map[string]bool{}
	for _, factory := range p.DataSources(context.Background()) {
		var metaResp datasource.MetadataResponse
		factory().Metadata(context.Background(), datasource.MetadataRequest{ProviderTypeName: "scaling"}, &metaResp)
		names[metaResp.TypeName] = true
	}
	return names
}

func TestInboundIntegrationDataSourceRegistered(t *testing.T) {
	t.Parallel()
	if !registeredDataSourceTypeNames(t)["scaling_inbound_integration"] {
		t.Errorf("scaling_inbound_integration data source not registered with the provider")
	}
}
