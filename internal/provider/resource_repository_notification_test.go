package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRepositoryNotification_webhook(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "notif-repo"
  visibility  = "private"
  description = ""
}

resource "quay_repository_notification" "test" {
  repository = quay_repository.test.id
  event      = "repo_push"
  method     = "webhook"
  title      = "push webhook"

  config = {
    url = "https://httpbin.org/post"
  }

  event_config = {}
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository_notification.test", "event", "repo_push"),
					resource.TestCheckResourceAttr("quay_repository_notification.test", "method", "webhook"),
					resource.TestCheckResourceAttr("quay_repository_notification.test", "title", "push webhook"),
					resource.TestCheckResourceAttrSet("quay_repository_notification.test", "uuid"),
				),
			},
		},
	})
}

func TestAccRepositoryNotification_vulnerability(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "vuln-repo"
  visibility  = "private"
  description = ""
}

resource "quay_repository_notification" "test" {
  repository = quay_repository.test.id
  event      = "vulnerability_found"
  method     = "webhook"
  title      = "vuln alert"

  config = {
    url = "https://httpbin.org/post"
  }

  event_config = {
    level = "High"
  }
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("quay_repository_notification.test", "event", "vulnerability_found"),
					resource.TestCheckResourceAttr("quay_repository_notification.test", "event_config.level", "High"),
					resource.TestCheckResourceAttrSet("quay_repository_notification.test", "uuid"),
				),
			},
		},
	})
}
