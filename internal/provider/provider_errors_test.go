package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestProvider_MissingToken verifies the provider returns a clear error when
// no token is supplied (neither in config nor environment).
func TestProvider_MissingToken(t *testing.T) {
	// Deliberately unset QUAY_TOKEN so the provider has nothing to use.
	t.Setenv("QUAY_TOKEN", "")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "quay" {
  url = "https://quay.io"
  # token deliberately omitted
}

data "quay_user" "current" {}
`,
				ExpectError: regexp.MustCompile(`Missing Quay API Token`),
			},
		},
	})
}

// TestProvider_InvalidToken verifies that an invalid token results in a
// meaningful error from the API (401) rather than a panic or empty state.
func TestProvider_InvalidToken(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
provider "quay" {
  url   = "` + testState.quayURL + `"
  token = "this-token-is-completely-invalid"
}

data "quay_user" "current" {}
`,
				ExpectError: regexp.MustCompile(`(?i)(401|Unauthorized|invalid_token|Error Reading User)`),
			},
		},
	})
}
