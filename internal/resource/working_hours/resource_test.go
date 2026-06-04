package working_hours_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

func TestAccWorkingHours_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkingHoursBasic("tf-acc-working-hours"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("scaling_working_hours.test", "id"),
					resource.TestCheckResourceAttr("scaling_working_hours.test", "name", "tf-acc-working-hours"),
					resource.TestCheckResourceAttr("scaling_working_hours.test", "timezone", "Europe/London"),
					resource.TestCheckResourceAttrSet("scaling_working_hours.test", "org_id"),
					resource.TestCheckResourceAttr("scaling_working_hours.test", "window.#", "1"),
					resource.TestCheckResourceAttr("scaling_working_hours.test", "window.0.start", "09:00"),
					resource.TestCheckResourceAttr("scaling_working_hours.test", "window.0.end", "17:00"),
					resource.TestCheckResourceAttr("scaling_working_hours.test", "window.0.days.#", "5"),
					resource.TestCheckResourceAttrSet("scaling_working_hours.test", "created_at"),
					resource.TestCheckResourceAttrSet("scaling_working_hours.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccWorkingHours_update(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkingHoursBasic("tf-acc-working-hours-update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaling_working_hours.test", "name", "tf-acc-working-hours-update"),
					resource.TestCheckResourceAttr("scaling_working_hours.test", "timezone", "Europe/London"),
				),
			},
			{
				Config: testAccWorkingHoursUpdated("tf-acc-working-hours-updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaling_working_hours.test", "name", "tf-acc-working-hours-updated"),
					resource.TestCheckResourceAttr("scaling_working_hours.test", "timezone", "America/New_York"),
					resource.TestCheckResourceAttr("scaling_working_hours.test", "window.#", "2"),
				),
			},
		},
	})
}

func testAccWorkingHoursBasic(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_working_hours" "test" {
  name     = %q
  timezone = "Europe/London"

  window {
    days  = [1, 2, 3, 4, 5]
    start = "09:00"
    end   = "17:00"
  }
}
`, testutil.ProviderConfig(), name)
}

func testAccWorkingHoursUpdated(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_working_hours" "test" {
  name     = %q
  timezone = "America/New_York"

  window {
    days  = [1, 2, 3, 4, 5]
    start = "08:00"
    end   = "18:00"
  }
  window {
    days  = [6, 7]
    start = "10:00"
    end   = "14:00"
  }
}
`, testutil.ProviderConfig(), name)
}
