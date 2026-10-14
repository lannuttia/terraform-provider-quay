package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepositoryTeamPermission_basic(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig() + `
resource "quay_organization" "test" {
  name  = "tf-acc-perm-org"
  email = "tf-acc-perm@example.com"
}

resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "perm-repo"
  visibility  = "private"
  description = ""
}

resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "readers"
  role     = "member"
}

resource "quay_repository_team_permission" "test" {
  repository = quay_repository.test.id
  team_name  = quay_organization_team.test.name
  role       = "read"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository_team_permission.test", "role", "read"),
				),
			},
			// Update role
			{
				Config: providerConfig() + `
resource "quay_organization" "test" {
  name  = "tf-acc-perm-org"
  email = "tf-acc-perm@example.com"
}

resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "perm-repo"
  visibility  = "private"
  description = ""
}

resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "readers"
  role     = "member"
}

resource "quay_repository_team_permission" "test" {
  repository = quay_repository.test.id
  team_name  = quay_organization_team.test.name
  role       = "write"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository_team_permission.test", "role", "write"),
				),
			},
		},
	})
}

func TestAccRepositoryUserPermission_basic(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_organization" "test" {
  name  = "tf-acc-userperm-org"
  email = "tf-acc-userperm@example.com"
}

resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "userperm-repo"
  visibility  = "private"
  description = ""
}

# Grant the admin user (quayadmin) explicit read access.
resource "quay_repository_user_permission" "test" {
  repository = quay_repository.test.id
  username   = "quayadmin"
  role       = "read"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository_user_permission.test", "role", "read"),
					resource.TestCheckResourceAttr("quay_repository_user_permission.test", "username", "quayadmin"),
				),
			},
		},
	})
}
