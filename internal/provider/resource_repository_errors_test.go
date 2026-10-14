package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

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

func TestAccRepository_DeletedOutOfBand(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository" "oob" {
  namespace   = "quayadmin"
  name        = "oob-repo"
  visibility  = "private"
  description = ""
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository.oob", "name", "oob-repo"),
				),
			},
			{
				RefreshState: true,
			},
		},
	})
}
