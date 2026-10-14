package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganization_basic(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig() + `
resource "quay_organization" "test" {
  name  = "tf-acc-org"
  email = "tf-acc@example.com"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization.test", "name", "tf-acc-org"),
					resource.TestCheckResourceAttr("quay_organization.test", "email", "tf-acc@example.com"),
				),
			},
			// Update
			{
				Config: providerConfig() + `
resource "quay_organization" "test" {
  name  = "tf-acc-org"
  email = "tf-acc-updated@example.com"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization.test", "email", "tf-acc-updated@example.com"),
				),
			},
		},
	})
}

func TestAccOrganization_datasource(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_organization" "test" {
  name  = "tf-acc-ds-org"
  email = "tf-acc-ds@example.com"
}

data "quay_organization" "test" {
  name = quay_organization.test.name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.quay_organization.test", "name", "tf-acc-ds-org"),
					resource.TestCheckResourceAttrSet("data.quay_organization.test", "is_organization"),
				),
			},
		},
	})
}
