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

func TestAccEscalationPolicy_workingHoursCondition(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyConditional("tf-acc-policy-cond"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaling_escalation_policy.test", "step.#", "2"),
					resource.TestCheckResourceAttr("scaling_escalation_policy.test", "step.0.condition.#", "1"),
					resource.TestCheckResourceAttrPair(
						"scaling_escalation_policy.test", "step.0.condition.0.working_hours_id",
						"scaling_working_hours.office", "id",
					),
					resource.TestCheckResourceAttr("scaling_escalation_policy.test", "step.0.condition.0.when", "during"),
					resource.TestCheckResourceAttr("scaling_escalation_policy.test", "step.1.condition.#", "0"),
				),
			},
			{
				// Re-applying the same config must be a no-op: conditions round-trip.
				Config:   testAccEscalationPolicyConditional("tf-acc-policy-cond"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccEscalationPolicy_conditionRemovalSurfacesDiff(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccEscalationPolicyConditional("tf-acc-policy-cond-rm"),
				Check: resource.TestCheckResourceAttr(
					"scaling_escalation_policy.test", "step.0.condition.#", "1",
				),
			},
			{
				// Dropping the condition block from config must surface a plan
				// (and clear it on apply) rather than silently preserving it.
				Config: testAccEscalationPolicyUnconditional("tf-acc-policy-cond-rm"),
				Check: resource.TestCheckResourceAttr(
					"scaling_escalation_policy.test", "step.0.condition.#", "0",
				),
			},
		},
	})
}

func testAccEscalationPolicyConditional(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_oncall_schedule" "dep" {
  name     = "%s-schedule"
  timezone = "UTC"
}

resource "scaling_working_hours" "office" {
  name     = "%s-hours"
  timezone = "Europe/London"

  window {
    days  = [1, 2, 3, 4, 5]
    start = "09:00"
    end   = "17:00"
  }
}

resource "scaling_escalation_policy" "test" {
  name = %q

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.dep.id
    escalate_after_seconds = 300

    condition {
      working_hours_id = scaling_working_hours.office.id
      when             = "during"
    }
  }

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.dep.id
    escalate_after_seconds = 600
  }
}
`, testutil.ProviderConfig(), name, name, name)
}

func testAccEscalationPolicyUnconditional(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_oncall_schedule" "dep" {
  name     = "%s-schedule"
  timezone = "UTC"
}

resource "scaling_working_hours" "office" {
  name     = "%s-hours"
  timezone = "Europe/London"

  window {
    days  = [1, 2, 3, 4, 5]
    start = "09:00"
    end   = "17:00"
  }
}

resource "scaling_escalation_policy" "test" {
  name = %q

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.dep.id
    escalate_after_seconds = 300
  }

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.dep.id
    escalate_after_seconds = 600
  }
}
`, testutil.ProviderConfig(), name, name, name)
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
