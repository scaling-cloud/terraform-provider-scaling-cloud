package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func registeredResourceTypeNames(t *testing.T) map[string]bool {
	t.Helper()
	p := New("test")()
	names := map[string]bool{}
	for _, factory := range p.Resources(context.Background()) {
		var metaResp resource.MetadataResponse
		factory().Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "scaling"}, &metaResp)
		names[metaResp.TypeName] = true
	}
	return names
}

func TestComponentResourceRegistered(t *testing.T) {
	t.Parallel()
	if !registeredResourceTypeNames(t)["scaling_component"] {
		t.Errorf("scaling_component not registered with the provider")
	}
}

func TestInboundIntegrationResourceRegistered(t *testing.T) {
	t.Parallel()
	if !registeredResourceTypeNames(t)["scaling_inbound_integration"] {
		t.Errorf("scaling_inbound_integration not registered with the provider")
	}
}
