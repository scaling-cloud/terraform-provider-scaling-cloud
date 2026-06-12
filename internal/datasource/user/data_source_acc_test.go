package user_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

// TestAccUserDataSource_basic resolves a user by email. Users are provisioned
// out-of-band (org membership), so the email is supplied via
// SCALING_CLOUD_ACC_USER_EMAIL and the test skips when it is unset.
func TestAccUserDataSource_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	email := os.Getenv("SCALING_CLOUD_ACC_USER_EMAIL")
	if email == "" {
		t.Skip("SCALING_CLOUD_ACC_USER_EMAIL not set, skipping happy-path lookup")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSource(email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.scaling_user.test", "id"),
					resource.TestCheckResourceAttr("data.scaling_user.test", "email", email),
				),
			},
		},
	})
}

func TestAccUserDataSource_notFound(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserDataSource("tf-acc-does-not-exist@example.invalid"),
				ExpectError: regexp.MustCompile(`no user found`),
			},
		},
	})
}

func testAccUserDataSource(email string) string {
	return fmt.Sprintf(`
%s

data "scaling_user" "test" {
  email = %q
}
`, testutil.ProviderConfig(), email)
}
