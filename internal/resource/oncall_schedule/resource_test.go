package oncall_schedule_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

func TestAccOncallSchedule_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOncallScheduleBasic("tf-acc-schedule"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("scaling_oncall_schedule.test", "id"),
					resource.TestCheckResourceAttr("scaling_oncall_schedule.test", "name", "tf-acc-schedule"),
					resource.TestCheckResourceAttr("scaling_oncall_schedule.test", "timezone", "America/New_York"),
					resource.TestCheckResourceAttrSet("scaling_oncall_schedule.test", "org_id"),
					resource.TestCheckResourceAttrSet("scaling_oncall_schedule.test", "created_at"),
					resource.TestCheckResourceAttrSet("scaling_oncall_schedule.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccOncallSchedule_update(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOncallScheduleBasic("tf-acc-schedule-update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaling_oncall_schedule.test", "name", "tf-acc-schedule-update"),
					resource.TestCheckResourceAttr("scaling_oncall_schedule.test", "timezone", "America/New_York"),
				),
			},
			{
				Config: testAccOncallScheduleUpdated("tf-acc-schedule-updated", "Europe/London"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaling_oncall_schedule.test", "name", "tf-acc-schedule-updated"),
					resource.TestCheckResourceAttr("scaling_oncall_schedule.test", "timezone", "Europe/London"),
				),
			},
		},
	})
}

func testAccOncallScheduleBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_oncall_schedule" "test" {
  name     = %q
  timezone = "America/New_York"
}
`, testutil.ProviderConfig(), name)
}

func testAccOncallScheduleUpdated(name, timezone string) string {
	return fmt.Sprintf(`
%s

resource "scaling_oncall_schedule" "test" {
  name        = %q
  timezone    = %q
  description = "Updated description"
}
`, testutil.ProviderConfig(), name, timezone)
}
