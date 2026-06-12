package working_hours_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

func TestAccWorkingHoursDataSource_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccWorkingHoursDataSource("tf-acc-wh-ds"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaling_working_hours.test", "id",
						"scaling_working_hours.test", "id",
					),
					resource.TestCheckResourceAttr("data.scaling_working_hours.test", "name", "tf-acc-wh-ds"),
					resource.TestCheckResourceAttr("data.scaling_working_hours.test", "timezone", "Europe/London"),
				),
			},
		},
	})
}

func TestAccWorkingHoursDataSource_notFound(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

data "scaling_working_hours" "test" {
  name = "tf-acc-does-not-exist-zzzzzz"
}
`, testutil.ProviderConfig()),
				ExpectError: regexp.MustCompile(`no working-hours set found`),
			},
		},
	})
}

func testAccWorkingHoursDataSource(name string) string {
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

data "scaling_working_hours" "test" {
  name = scaling_working_hours.test.name
}
`, testutil.ProviderConfig(), name)
}
