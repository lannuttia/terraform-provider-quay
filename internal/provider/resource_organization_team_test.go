package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganizationTeam_basic(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-team-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc-team@example.com"
}

resource "quay_organization_team" "test" {
  org_name    = quay_organization.test.name
  name        = "developers"
  role        = "member"
  description = "Developer team"
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization_team.test", "name", "developers"),
					resource.TestCheckResourceAttr("quay_organization_team.test", "role", "member"),
					resource.TestCheckResourceAttr("quay_organization_team.test", "id", orgName+"/developers"),
				),
			},
			// Update role
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc-team@example.com"
}

resource "quay_organization_team" "test" {
  org_name    = quay_organization.test.name
  name        = "developers"
  role        = "creator"
  description = "Developer team"
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization_team.test", "role", "creator"),
				),
			},
		},
	})
}
