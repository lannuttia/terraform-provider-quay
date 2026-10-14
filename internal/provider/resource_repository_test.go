package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepository_basic(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-repo-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc-repo@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "my-repo"
  visibility  = "private"
  description = "acceptance test repo"
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository.test", "name", "my-repo"),
					resource.TestCheckResourceAttr("quay_repository.test", "visibility", "private"),
					resource.TestCheckResourceAttr("quay_repository.test", "description", "acceptance test repo"),
					resource.TestCheckResourceAttr("quay_repository.test", "id", orgName+"/my-repo"),
				),
			},
			// Update description
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc-repo@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "my-repo"
  visibility  = "private"
  description = "updated description"
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository.test", "description", "updated description"),
				),
			},
			// Change visibility
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc-repo@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "my-repo"
  visibility  = "public"
  description = "updated description"
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository.test", "visibility", "public"),
				),
			},
		},
	})
}

func TestAccRepository_datasource(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-ds-repo-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc-ds-repo@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "ds-repo"
  visibility  = "private"
  description = "datasource test"
}
data "quay_repository" "test" {
  namespace = quay_repository.test.namespace
  name      = quay_repository.test.name
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.quay_repository.test", "visibility", "private"),
					resource.TestCheckResourceAttr("data.quay_repository.test", "description", "datasource test"),
				),
			},
		},
	})
}
