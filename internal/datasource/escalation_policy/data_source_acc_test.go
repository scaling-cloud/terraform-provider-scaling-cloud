package escalation_policy_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

func TestAccEscalationPolicyDataSource_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyDataSource("tf-acc-ep-ds"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaling_escalation_policy.test", "id",
						"scaling_escalation_policy.test", "id",
					),
					resource.TestCheckResourceAttr("data.scaling_escalation_policy.test", "name", "tf-acc-ep-ds"),
					resource.TestCheckResourceAttrSet("data.scaling_escalation_policy.test", "org_id"),
				),
			},
		},
	})
}

func TestAccEscalationPolicyDataSource_notFound(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

data "scaling_escalation_policy" "test" {
  name = "tf-acc-does-not-exist-zzzzzz"
}
`, testutil.ProviderConfig()),
				ExpectError: regexp.MustCompile(`no escalation policy found`),
			},
		},
	})
}

func testAccEscalationPolicyDataSource(name string) string {
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

data "scaling_escalation_policy" "test" {
  name = scaling_escalation_policy.test.name
}
`, testutil.ProviderConfig(), name, name)
}
