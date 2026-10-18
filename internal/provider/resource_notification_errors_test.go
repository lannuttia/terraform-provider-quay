package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepositoryNotification_NoMethodBlock(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository_notification" "bad" {
  repository = "myorg/myrepo"
  repo_push  = {}
}
`,
				ExpectError: regexp.MustCompile(`(?i)(exactly one of)`),
			},
		},
	})
}

func TestAccRepositoryNotification_MultipleMethodBlocks(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository_notification" "bad" {
  repository = "myorg/myrepo"
  webhook    = { url = "https://example.com" }
  slack      = { url = "https://hooks.slack.com/services/T/B/X" }
  repo_push  = {}
}
`,
				ExpectError: regexp.MustCompile(`(?i)(exactly one of)`),
			},
		},
	})
}

func TestAccRepositoryNotification_NoEventBlock(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository_notification" "bad" {
  repository = "myorg/myrepo"
  webhook    = { url = "https://example.com" }
}
`,
				ExpectError: regexp.MustCompile(`(?i)(exactly one of)`),
			},
		},
	})
}

func TestAccRepositoryNotification_MultipleEventBlocks(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository_notification" "bad" {
  repository    = "myorg/myrepo"
  webhook       = { url = "https://example.com" }
  repo_push     = {}
  build_success = {}
}
`,
				ExpectError: regexp.MustCompile(`(?i)(exactly one of)`),
			},
		},
	})
}

func TestAccRepositoryNotification_InvalidVulnerabilityLevel(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository_notification" "bad" {
  repository          = "myorg/myrepo"
  webhook             = { url = "https://example.com" }
  vulnerability_found = { level = "Catastrophic" }
}
`,
				ExpectError: regexp.MustCompile(`(?i)(invalid|value must be one of)`),
			},
		},
	})
}

func TestAccRepositoryNotification_RepositoryNotFound(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository_notification" "bad" {
  repository = "this-org-does-not-exist-tf-acc/also-missing"
  webhook    = { url = "https://example.com" }
  repo_push  = {}
}
`,
				ExpectError: regexp.MustCompile(`(?i)(not found|404|Error Creating Repository Notification)`),
			},
		},
	})
}
