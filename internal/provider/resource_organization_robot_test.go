package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationRobot_basic(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig() + `
resource "quay_organization" "test" {
  name  = "tf-acc-robot-org"
  email = "tf-acc-robot@example.com"
}

resource "quay_organization_robot" "test" {
  org_name    = quay_organization.test.name
  short_name  = "ci-bot"
  description = "CI robot account"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization_robot.test", "short_name", "ci-bot"),
					resource.TestCheckResourceAttr("quay_organization_robot.test", "description", "CI robot account"),
					resource.TestCheckResourceAttr("quay_organization_robot.test", "full_name", "tf-acc-robot-org+ci-bot"),
					resource.TestCheckResourceAttrSet("quay_organization_robot.test", "token"),
				),
			},
		},
	})
}
