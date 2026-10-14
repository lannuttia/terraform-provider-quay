package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// notificationImportID builds the import ID (namespace/repo_name/uuid) from state.
func notificationImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		r, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %q not found in state", resourceName)
		}
		return fmt.Sprintf("%s/%s", r.Primary.Attributes["repository"], r.Primary.Attributes["uuid"]), nil
	}
}

// ---------------------------------------------------------------------------
// quay_repository
// ---------------------------------------------------------------------------

func TestAccRepository_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	cfg := providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "import-repo"
  visibility  = "private"
  description = "import test"
}`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				Config:            cfg,
				ResourceName:      "quay_repository.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRepository_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	cfg := providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "import-repo"
  visibility  = "private"
  description = "import test"
}`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				Config:          cfg,
				ResourceName:    "quay_repository.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
			},
		},
	})
}

func TestAccRepository_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	cfg := providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "import-repo"
  visibility  = "private"
  description = "import test"
}`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				Config:          cfg,
				ResourceName:    "quay_repository.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// quay_repository_notification
// ---------------------------------------------------------------------------

func TestAccRepositoryNotification_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	cfg := providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "import-notif-repo"
  visibility  = "private"
  description = ""
}
resource "quay_repository_notification" "test" {
  repository   = quay_repository.test.id
  event        = "repo_push"
  method       = "webhook"
  title        = "import test"
  config       = { url = "https://httpbin.org/post" }
  event_config = {}
}`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				Config:            cfg,
				ResourceName:      "quay_repository_notification.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: notificationImportID("quay_repository_notification.test"),
			},
		},
	})
}

func TestAccRepositoryNotification_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	cfg := providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "import-notif-repo"
  visibility  = "private"
  description = ""
}
resource "quay_repository_notification" "test" {
  repository   = quay_repository.test.id
  event        = "repo_push"
  method       = "webhook"
  title        = "import test"
  config       = { url = "https://httpbin.org/post" }
  event_config = {}
}`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				Config:            cfg,
				ResourceName:      "quay_repository_notification.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateIdFunc: notificationImportID("quay_repository_notification.test"),
			},
		},
	})
}

func TestAccRepositoryNotification_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	cfg := providerConfig() + `
resource "quay_repository" "test" {
  namespace   = "quayadmin"
  name        = "import-notif-repo"
  visibility  = "private"
  description = ""
}
resource "quay_repository_notification" "test" {
  repository   = quay_repository.test.id
  event        = "repo_push"
  method       = "webhook"
  title        = "import test"
  config       = { url = "https://httpbin.org/post" }
  event_config = {}
}`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				Config:            cfg,
				ResourceName:      "quay_repository_notification.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				GenerateConfig:    true,
				ImportStateIdFunc: notificationImportID("quay_repository_notification.test"),
			},
		},
	})
}
