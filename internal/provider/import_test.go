package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOrganization_import(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
`, orgName),
			},
			{
				ResourceName:      "quay_organization.test",
				ImportState:       true,
				ImportStateVerify: true,
				// tag_expiration_s may differ between create response and import read.
				ImportStateVerifyIgnore: []string{"tag_expiration_s", "invoice_email", "invoice_email_address", "contact_email"},
			},
		},
	})
}

func TestAccRepository_import(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-repo-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "import-repo"
  visibility  = "private"
  description = "import test"
}
`, orgName),
			},
			{
				ResourceName:            "quay_repository.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdPrefix:     orgName + "/",
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAccOrganizationRobot_import(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-robot-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_robot" "test" {
  org_name    = quay_organization.test.name
  short_name  = "import-bot"
  description = "import test robot"
}
`, orgName),
			},
			{
				ResourceName:      "quay_organization_robot.test",
				ImportState:       true,
				ImportStateVerify: true,
				// token is write-only — never returned on read after creation.
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func TestAccOrganizationTeam_import(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-team-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_team" "test" {
  org_name    = quay_organization.test.name
  name        = "import-team"
  role        = "member"
  description = "import test team"
}
`, orgName),
			},
			{
				ResourceName:      "quay_organization_team.test",
				ImportState:       true,
				ImportStateVerify: true,
				// The Quay API has no GET single team — role and description
				// are not returned on read, so we can't verify them on import.
				ImportStateVerifyIgnore: []string{"role", "description"},
			},
		},
	})
}

// testAccNotificationImportID builds the import ID (namespace/repo/uuid) from state.
func testAccNotificationImportID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		r, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %q not found in state", resourceName)
		}
		repo := r.Primary.Attributes["repository"]
		uuid := r.Primary.Attributes["uuid"]
		return fmt.Sprintf("%s/%s", repo, uuid), nil
	}
}

func TestAccRepositoryNotification_import(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-notif-org")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				Config: providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "import-notif-repo"
  visibility  = "private"
  description = ""
}
resource "quay_repository_notification" "test" {
  repository = quay_repository.test.id
  event      = "repo_push"
  method     = "webhook"
  title      = "import test"
  config     = { url = "https://httpbin.org/post" }
  event_config = {}
}
`, orgName),
			},
			{
				ResourceName:      "quay_repository_notification.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Import ID is namespace/repo_name/uuid — derived from the
				// resource's id attribute (repo) plus uuid split into 3 parts.
				ImportStateIdFunc: testAccNotificationImportID("quay_repository_notification.test"),
			},
		},
	})
}
