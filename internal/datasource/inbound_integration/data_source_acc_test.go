package inbound_integration_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

// TestAccInboundIntegrationDataSource_basic resolves a pre-provisioned inbound
// integration by name. Inbound integrations have no create API (they are
// installed out-of-band), so the integration's name is supplied via
// SCALING_CLOUD_ACC_INBOUND_INTEGRATION_NAME and the test skips when it is unset.
func TestAccInboundIntegrationDataSource_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	name := os.Getenv("SCALING_CLOUD_ACC_INBOUND_INTEGRATION_NAME")
	if name == "" {
		t.Skip("SCALING_CLOUD_ACC_INBOUND_INTEGRATION_NAME not set, skipping happy-path lookup")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInboundIntegrationDataSource(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaling_inbound_integration.test", "id"),
					resource.TestCheckResourceAttr("data.scaling_inbound_integration.test", "name", name),
					resource.TestCheckResourceAttrSet("data.scaling_inbound_integration.test", "org_id"),
					resource.TestCheckResourceAttrSet("data.scaling_inbound_integration.test", "component_id"),
				),
			},
		},
	})
}

// TestAccInboundIntegrationDataSource_notFound asserts that looking up a name
// with no match fails loudly rather than returning empty state.
func TestAccInboundIntegrationDataSource_notFound(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccInboundIntegrationDataSource("tf-acc-does-not-exist-zzzzzz"),
				ExpectError: regexp.MustCompile(`no inbound integration found`),
			},
		},
	})
}

func testAccInboundIntegrationDataSource(name string) string {
	return fmt.Sprintf(`
%s

data "scaling_inbound_integration" "test" {
  name = %q
}
`, testutil.ProviderConfig(), name)
}
