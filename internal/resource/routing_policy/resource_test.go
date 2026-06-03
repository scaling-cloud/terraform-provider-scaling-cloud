package routing_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

func TestAccRoutingPolicy_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingPolicyBasic("tf-acc-routing"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("scaling_routing_policy.test", "id"),
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "name", "tf-acc-routing"),
					resource.TestCheckResourceAttrSet("scaling_routing_policy.test", "org_id"),
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "is_default", "false"),
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "rule.#", "4"),
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "rule.0.severity", "critical"),
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "rule.0.outcome", "incident"),
					resource.TestCheckResourceAttrSet("scaling_routing_policy.test", "created_at"),
					resource.TestCheckResourceAttrSet("scaling_routing_policy.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccRoutingPolicy_update(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingPolicyBasic("tf-acc-routing-update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "name", "tf-acc-routing-update"),
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "rule.0.outcome", "incident"),
				),
			},
			{
				Config: testAccRoutingPolicyUpdated("tf-acc-routing-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "name", "tf-acc-routing-updated"),
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "description", "Updated routing"),
					resource.TestCheckResourceAttr("scaling_routing_policy.test", "rule.0.outcome", "provisional_page"),
				),
			},
		},
	})
}

func testAccRoutingPolicyBasic(name string) string {
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
`, testutil.ProviderConfig(), name)
}

func testAccRoutingPolicyUpdated(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_routing_policy" "test" {
  name        = %q
  description = "Updated routing"

  rule {
    severity = "critical"
    outcome  = "provisional_page"
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
`, testutil.ProviderConfig(), name)
}
