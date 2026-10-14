package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganization_basic(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc@example.com"
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization.test", "name", orgName),
					resource.TestCheckResourceAttr("quay_organization.test", "email", "tf-acc@example.com"),
				),
			},
			// Update
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc-updated@example.com"
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization.test", "email", "tf-acc-updated@example.com"),
				),
			},
		},
	})
}

func TestAccOrganization_datasource(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-ds-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "tf-acc-ds@example.com"
}

data "quay_organization" "test" {
  name = quay_organization.test.name
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.quay_organization.test", "name", orgName),
					resource.TestCheckResourceAttrSet("data.quay_organization.test", "is_organization"),
				),
			},
		},
	})
}
