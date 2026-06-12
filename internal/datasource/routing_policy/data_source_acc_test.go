package routing_policy_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

func TestAccRoutingPolicyDataSource_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingPolicyDataSource("tf-acc-rp-ds"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaling_routing_policy.test", "id",
						"scaling_routing_policy.test", "id",
					),
					resource.TestCheckResourceAttr("data.scaling_routing_policy.test", "name", "tf-acc-rp-ds"),
					resource.TestCheckResourceAttr("data.scaling_routing_policy.test", "is_default", "false"),
				),
			},
		},
	})
}

func TestAccRoutingPolicyDataSource_notFound(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

data "scaling_routing_policy" "test" {
  name = "tf-acc-does-not-exist-zzzzzz"
}
`, testutil.ProviderConfig()),
				ExpectError: regexp.MustCompile(`no routing policy found`),
			},
		},
	})
}

func testAccRoutingPolicyDataSource(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_routing_policy" "test" {
  name = %q

  rule {
    severity = "critical"
    outcome  = "incident"
  }
  rule {
    severity = "high"
    outcome  = "provisional_page"
  }
  rule {
    severity = "medium"
    outcome  = "notification"
  }
  rule {
    severity = "low"
    outcome  = "drop"
  }
}

data "scaling_routing_policy" "test" {
  name = scaling_routing_policy.test.name
}
`, testutil.ProviderConfig(), name)
}
