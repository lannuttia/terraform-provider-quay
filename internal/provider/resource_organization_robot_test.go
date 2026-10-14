package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationRobot_basic(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-robot-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc-robot@example.com"
}

resource "quay_organization_robot" "test" {
  org_name    = quay_organization.test.name
  short_name  = "ci-bot"
  description = "CI robot account"
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization_robot.test", "short_name", "ci-bot"),
					resource.TestCheckResourceAttr("quay_organization_robot.test", "description", "CI robot account"),
					resource.TestCheckResourceAttr("quay_organization_robot.test", "full_name", orgName+"+ci-bot"),
					resource.TestCheckResourceAttrSet("quay_organization_robot.test", "token"),
				),
			},
		},
	})
}
