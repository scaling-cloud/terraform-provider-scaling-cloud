package oncall_schedule_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

func TestAccOncallScheduleDataSource_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOncallScheduleDataSource("tf-acc-sch-ds"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaling_oncall_schedule.test", "id",
						"scaling_oncall_schedule.test", "id",
					),
					resource.TestCheckResourceAttr("data.scaling_oncall_schedule.test", "name", "tf-acc-sch-ds"),
					resource.TestCheckResourceAttr("data.scaling_oncall_schedule.test", "timezone", "America/New_York"),
				),
			},
		},
	})
}

func TestAccOncallScheduleDataSource_notFound(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

data "scaling_oncall_schedule" "test" {
  name = "tf-acc-does-not-exist-zzzzzz"
}
`, testutil.ProviderConfig()),
				ExpectError: regexp.MustCompile(`no on-call schedule found`),
			},
		},
	})
}

func testAccOncallScheduleDataSource(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_oncall_schedule" "test" {
  name     = %q
  timezone = "America/New_York"
}

data "scaling_oncall_schedule" "test" {
  name = scaling_oncall_schedule.test.name
}
`, testutil.ProviderConfig(), name)
}
