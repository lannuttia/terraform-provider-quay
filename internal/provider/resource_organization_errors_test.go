package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrganization_DuplicateName(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-dup-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create the first instance.
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "first" {
  name  = %q
  email = "first@example.com"
}
`, orgName),
			},
			// Attempt to create a second resource with the same name — should fail.
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "first" {
  name  = %q
  email = "first@example.com"
}

resource "quay_organization" "second" {
  name  = %q
  email = "second@example.com"
}
`, orgName, orgName),
				ExpectError: regexp.MustCompile(`(?i)(already exists|conflict|409|Error Creating Organization)`),
			},
		},
	})
}

func TestAccOrganization_DataSourceNotFound(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
data "quay_organization" "missing" {
  name = "this-org-does-not-exist-tf-acc"
}
`,
				ExpectError: regexp.MustCompile(`(?i)(not found|404|Error Reading Organization)`),
			},
		},
	})
}

func TestAccOrganization_DeletedOutOfBand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-oob-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "oob" {
  name  = %q
  email = "oob@example.com"
}
`, orgName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization.oob", "name", orgName),
				),
			},
			{
				RefreshState: true,
			},
		},
	})
}
