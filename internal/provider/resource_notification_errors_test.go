package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepositoryNotification_InvalidEvent(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository_notification" "bad" {
  repository = "myorg/myrepo"
  event      = "not_a_real_event"
  method     = "webhook"
  config     = { url = "https://example.com" }
  event_config = {}
}
`,
				ExpectError: regexp.MustCompile(`(?i)(invalid|value must be one of)`),
			},
		},
	})
}

func TestAccRepositoryNotification_InvalidMethod(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository_notification" "bad" {
  repository = "myorg/myrepo"
  event      = "repo_push"
  method     = "carrier_pigeon"
  config     = { url = "https://example.com" }
  event_config = {}
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
  event      = "repo_push"
  method     = "webhook"
  config     = { url = "https://example.com" }
  event_config = {}
}
`,
				ExpectError: regexp.MustCompile(`(?i)(not found|404|Error Creating Repository Notification)`),
			},
		},
	})
}
