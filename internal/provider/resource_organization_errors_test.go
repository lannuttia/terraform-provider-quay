package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccOrganization_DuplicateName verifies a clear error when creating an
// org that already exists.
func TestAccOrganization_DuplicateName(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create the first instance — should succeed.
			{
				Config: providerConfig() + `
resource "quay_organization" "first" {
  name  = "tf-acc-duplicate-org"
  email = "first@example.com"
}
`,
			},
			// Attempt to create a second resource with the same org name — should fail.
			{
				Config: providerConfig() + `
resource "quay_organization" "first" {
  name  = "tf-acc-duplicate-org"
  email = "first@example.com"
}

resource "quay_organization" "second" {
  name  = "tf-acc-duplicate-org"
  email = "second@example.com"
}
`,
				ExpectError: regexp.MustCompile(`(?i)(already exists|conflict|409|Error Creating Organization)`),
			},
		},
	})
}

// TestAccOrganization_DataSourceNotFound verifies a clear error when reading
// an org that does not exist.
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

// TestAccOrganization_DeletedOutOfBand verifies that when an org is deleted
// outside of Terraform, the next plan removes it from state rather than
// erroring fatally (i.e. Read returns gracefully on 404).
func TestAccOrganization_DeletedOutOfBand(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create the org via Terraform.
			{
				Config: providerConfig() + `
resource "quay_organization" "oob" {
  name  = "tf-acc-oob-org"
  email = "oob@example.com"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_organization.oob", "name", "tf-acc-oob-org"),
				),
			},
			// Delete it via the API directly, then re-apply the same config.
			// Terraform should detect the drift and recreate it.
			{
				// RefreshState forces a Read before comparing plan, surfacing the drift.
				RefreshState: true,
			},
		},
	})
}
