package component_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/testutil"
)

func TestAccComponentDataSource_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComponentDataSource("tf-acc-cmp-ds"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaling_component.test", "id",
						"scaling_component.test", "id",
					),
					resource.TestCheckResourceAttr("data.scaling_component.test", "name", "tf-acc-cmp-ds"),
					resource.TestCheckResourceAttrSet("data.scaling_component.test", "operational_status"),
					resource.TestCheckResourceAttr("data.scaling_component.test", "aliases.#", "1"),
				),
			},
		},
	})
}

func TestAccComponentDataSource_byAlias(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccComponentDataSourceByAlias("tf-acc-cmp-alias"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.scaling_component.by_alias", "id",
						"scaling_component.test", "id",
					),
					resource.TestCheckResourceAttr("data.scaling_component.by_alias", "name", "tf-acc-cmp-alias"),
					resource.TestCheckResourceAttr("data.scaling_component.by_alias", "aliases.#", "1"),
				),
			},
		},
	})
}

func TestAccComponentDataSource_notFound(t *testing.T) {
	testutil.TestAccPreCheck(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
%s

data "scaling_component" "test" {
  name = "tf-acc-does-not-exist-zzzzzz"
}
`, testutil.ProviderConfig()),
				ExpectError: regexp.MustCompile(`no component found`),
			},
		},
	})
}

func testAccComponentDataSource(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_component" "test" {
  name    = %q
  aliases = ["%s-alias"]
}

data "scaling_component" "test" {
  name = scaling_component.test.name
}
`, testutil.ProviderConfig(), name, name)
}

func testAccComponentDataSourceByAlias(name string) string {
	return fmt.Sprintf(`
%s

resource "scaling_component" "test" {
  name    = %q
  aliases = ["%s-alias"]
}

data "scaling_component" "by_alias" {
  name  = scaling_component.test.name
  alias = "%s-alias"
}
`, testutil.ProviderConfig(), name, name, name)
}
