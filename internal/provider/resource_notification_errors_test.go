package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccRepositoryNotification_InvalidEvent verifies that an unrecognised
// event type is caught at plan time by the schema validator.
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

// TestAccRepositoryNotification_InvalidMethod verifies that an unrecognised
// method is caught at plan time by the schema validator.
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

// TestAccRepositoryNotification_RepositoryNotFound verifies a clear error when
// creating a notification on a repository that does not exist.
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
