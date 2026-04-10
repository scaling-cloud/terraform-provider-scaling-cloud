package escalation_policy_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

func TestAccEscalationPolicy_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyBasic("tf-acc-policy"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("scaling_escalation_policy.test", "id"),
					resource.TestCheckResourceAttr("scaling_escalation_policy.test", "name", "tf-acc-policy"),
					resource.TestCheckResourceAttrSet("scaling_escalation_policy.test", "org_id"),
					resource.TestCheckResourceAttrSet("scaling_escalation_policy.test", "created_at"),
					resource.TestCheckResourceAttrSet("scaling_escalation_policy.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccEscalationPolicy_update(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyBasic("tf-acc-policy-update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaling_escalation_policy.test", "name", "tf-acc-policy-update"),
				),
			},
			{
				Config: testAccEscalationPolicyUpdated("tf-acc-policy-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaling_escalation_policy.test", "name", "tf-acc-policy-updated"),
					resource.TestCheckResourceAttr("scaling_escalation_policy.test", "description", "Updated policy"),
				),
			},
		},
	})
}

func testAccEscalationPolicyBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_oncall_schedule" "dep" {
  name     = "%s-schedule"
  timezone = "UTC"
}

resource "scaling_escalation_policy" "test" {
  name = %q

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.dep.id
    escalate_after_seconds = 300
  }
}
`, testutil.ProviderConfig(), name, name)
}

func testAccEscalationPolicyUpdated(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_oncall_schedule" "dep" {
  name     = "%s-schedule"
  timezone = "UTC"
}

resource "scaling_escalation_policy" "test" {
  name        = %q
  description = "Updated policy"

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.dep.id
    escalate_after_seconds = 600
  }
}
`, testutil.ProviderConfig(), name, name)
}
