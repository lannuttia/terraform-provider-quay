package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// All repository tests use "quayadmin" as the namespace — that user is always
// present in the testcontainer and doesn't require a quay_organization resource.

func TestAccRepository_basic(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "my-repo"
  visibility  = "private"
  description = "acceptance test repo"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository.test", "name", "my-repo"),
					resource.TestCheckResourceAttr("quay_repository.test", "visibility", "private"),
					resource.TestCheckResourceAttr("quay_repository.test", "description", "acceptance test repo"),
					resource.TestCheckResourceAttr("quay_repository.test", "id", "quayadmin/my-repo"),
				),
			},
			// Update description
			{
				Config: providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "my-repo"
  visibility  = "private"
  description = "updated description"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository.test", "description", "updated description"),
				),
			},
			// Change visibility
			{
				Config: providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "my-repo"
  visibility  = "public"
  description = "updated description"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository.test", "visibility", "public"),
				),
			},
		},
	})
}

func TestAccRepository_datasource(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "ds-repo"
  visibility  = "private"
  description = "datasource test"
}

data "quay_repository" "test" {
  namespace = quay_repository.test.namespace
  name      = quay_repository.test.name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.quay_repository.test", "visibility", "private"),
					resource.TestCheckResourceAttr("data.quay_repository.test", "description", "datasource test"),
				),
			},
		},
	})
}
