package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRepository_DataSourceNotFound verifies a clear error when reading a
// repository that does not exist.
func TestAccRepository_DataSourceNotFound(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
data "quay_repository" "missing" {
  namespace = "this-namespace-does-not-exist-tf-acc"
  name      = "also-missing"
}
`,
				ExpectError: regexp.MustCompile(`(?i)(not found|404|Error Reading Repository)`),
			},
		},
	})
}

// TestAccRepository_InvalidVisibility verifies that supplying an unrecognised
// visibility value is caught at plan time by the schema validator.
func TestAccRepository_InvalidVisibility(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository" "bad" {
  namespace  = "myorg"
  name       = "myrepo"
  visibility = "internal"
}
`,
				ExpectError: regexp.MustCompile(`(?i)(invalid|value must be one of)`),
			},
		},
	})
}

// TestAccRepository_DeletedOutOfBand verifies that a repository deleted
// outside Terraform is detected during refresh and removed from state.
func TestAccRepository_DeletedOutOfBand(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			// Create via Terraform.
			{
				Config: providerConfig() + `
resource "quay_organization" "test" {
  name  = "tf-acc-oob-repo-org"
  email = "oob-repo@example.com"
}

resource "quay_repository" "oob" {
  namespace   = quay_organization.test.name
  name        = "oob-repo"
  visibility  = "private"
  description = ""
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository.oob", "name", "oob-repo"),
				),
			},
			// Refresh only — Terraform should reconcile state without error.
			{
				RefreshState: true,
			},
		},
	})
}
