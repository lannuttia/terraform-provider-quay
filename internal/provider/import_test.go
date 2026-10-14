package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ---------------------------------------------------------------------------
// Helper: build an import ID for quay_repository_notification from state.
// The format is namespace/repo_name/uuid (3 slash-separated parts).
// ---------------------------------------------------------------------------

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
// quay_organization
// ---------------------------------------------------------------------------

func TestAccOrganization_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-org-cmd")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:            "quay_organization.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tag_expiration_s", "invoice_email", "invoice_email_address", "contact_email"},
			},
		},
	})
}

func TestAccOrganization_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-org-blk")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:            "quay_organization.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportBlockWithID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tag_expiration_s", "invoice_email", "invoice_email_address", "contact_email"},
			},
		},
	})
}

func TestAccOrganization_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-org-gen")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:    "quay_organization.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// quay_repository
// ---------------------------------------------------------------------------

func TestAccRepository_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-repo-cmd")
	cfg := providerConfig() + fmt.Sprintf(`
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
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_repository.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRepository_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-repo-blk")
	cfg := providerConfig() + fmt.Sprintf(`
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
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_repository.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRepository_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-repo-gen")
	cfg := providerConfig() + fmt.Sprintf(`
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
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:    "quay_repository.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// quay_organization_robot
// ---------------------------------------------------------------------------

func TestAccOrganizationRobot_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-robot-cmd")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_robot" "test" {
  org_name    = quay_organization.test.name
  short_name  = "import-bot"
  description = "import test robot"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:            "quay_organization_robot.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func TestAccOrganizationRobot_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-robot-blk")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_robot" "test" {
  org_name    = quay_organization.test.name
  short_name  = "import-bot"
  description = "import test robot"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:            "quay_organization_robot.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportBlockWithID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func TestAccOrganizationRobot_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-robot-gen")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_robot" "test" {
  org_name    = quay_organization.test.name
  short_name  = "import-bot"
  description = "import test robot"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:    "quay_organization_robot.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// quay_organization_team
// ---------------------------------------------------------------------------

func TestAccOrganizationTeam_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-team-cmd")
	cfg := providerConfig() + fmt.Sprintf(`
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
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:            "quay_organization_team.test",
				ImportState:             true,
				ImportStateVerify:       true,
				// Quay has no GET single team endpoint — role and description
				// cannot be read back after import.
				ImportStateVerifyIgnore: []string{"role", "description"},
			},
		},
	})
}

func TestAccOrganizationTeam_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-team-blk")
	cfg := providerConfig() + fmt.Sprintf(`
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
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:            "quay_organization_team.test",
				ImportState:             true,
				ImportStateKind:         resource.ImportBlockWithID,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"role", "description"},
			},
		},
	})
}

func TestAccOrganizationTeam_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-team-gen")
	cfg := providerConfig() + fmt.Sprintf(`
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
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:    "quay_organization_team.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// quay_organization_team_member
// ---------------------------------------------------------------------------

func TestAccOrganizationTeamMember_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-member-cmd")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "import-team"
  role     = "member"
}
resource "quay_organization_robot" "test" {
  org_name   = quay_organization.test.name
  short_name = "import-bot"
}
resource "quay_organization_team_member" "test" {
  org_name    = quay_organization.test.name
  team_name   = quay_organization_team.test.name
  member_name = quay_organization_robot.test.full_name
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_organization_team_member.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrganizationTeamMember_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-member-blk")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "import-team"
  role     = "member"
}
resource "quay_organization_robot" "test" {
  org_name   = quay_organization.test.name
  short_name = "import-bot"
}
resource "quay_organization_team_member" "test" {
  org_name    = quay_organization.test.name
  team_name   = quay_organization_team.test.name
  member_name = quay_organization_robot.test.full_name
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_organization_team_member.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrganizationTeamMember_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-member-gen")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "import-team"
  role     = "member"
}
resource "quay_organization_robot" "test" {
  org_name   = quay_organization.test.name
  short_name = "import-bot"
}
resource "quay_organization_team_member" "test" {
  org_name    = quay_organization.test.name
  team_name   = quay_organization_team.test.name
  member_name = quay_organization_robot.test.full_name
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:    "quay_organization_team_member.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// quay_organization_prototype_permission
// ---------------------------------------------------------------------------

func TestAccOrganizationPrototypePermission_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-proto-cmd")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "import-team"
  role     = "member"
}
resource "quay_organization_prototype_permission" "test" {
  org_name      = quay_organization.test.name
  role          = "read"
  delegate_kind = "team"
  delegate_name = quay_organization_team.test.name
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_organization_prototype_permission.test",
				ImportState:       true,
				ImportStateVerify: true,
				// Import ID is org_name/prototype_id — derive it from state.
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["quay_organization_prototype_permission.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return fmt.Sprintf("%s/%s", r.Primary.Attributes["org_name"], r.Primary.Attributes["id"]), nil
				},
			},
		},
	})
}

func TestAccOrganizationPrototypePermission_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-proto-blk")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "import-team"
  role     = "member"
}
resource "quay_organization_prototype_permission" "test" {
  org_name      = quay_organization.test.name
  role          = "read"
  delegate_kind = "team"
  delegate_name = quay_organization_team.test.name
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_organization_prototype_permission.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["quay_organization_prototype_permission.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return fmt.Sprintf("%s/%s", r.Primary.Attributes["org_name"], r.Primary.Attributes["id"]), nil
				},
			},
		},
	})
}

func TestAccOrganizationPrototypePermission_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-proto-gen")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "import-team"
  role     = "member"
}
resource "quay_organization_prototype_permission" "test" {
  org_name      = quay_organization.test.name
  role          = "read"
  delegate_kind = "team"
  delegate_name = quay_organization_team.test.name
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:    "quay_organization_prototype_permission.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					r, ok := s.RootModule().Resources["quay_organization_prototype_permission.test"]
					if !ok {
						return "", fmt.Errorf("resource not found in state")
					}
					return fmt.Sprintf("%s/%s", r.Primary.Attributes["org_name"], r.Primary.Attributes["id"]), nil
				},
			},
		},
	})
}

// ---------------------------------------------------------------------------
// quay_repository_notification
// ---------------------------------------------------------------------------

func TestAccRepositoryNotification_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-notif-cmd")
	cfg := providerConfig() + fmt.Sprintf(`
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
  repository   = quay_repository.test.id
  event        = "repo_push"
  method       = "webhook"
  title        = "import test"
  config       = { url = "https://httpbin.org/post" }
  event_config = {}
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
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
	orgName := rname("tf-acc-import-notif-blk")
	cfg := providerConfig() + fmt.Sprintf(`
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
  repository   = quay_repository.test.id
  event        = "repo_push"
  method       = "webhook"
  title        = "import test"
  config       = { url = "https://httpbin.org/post" }
  event_config = {}
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_repository_notification.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateVerify: true,
				ImportStateIdFunc: notificationImportID("quay_repository_notification.test"),
			},
		},
	})
}

func TestAccRepositoryNotification_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-notif-gen")
	cfg := providerConfig() + fmt.Sprintf(`
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
  repository   = quay_repository.test.id
  event        = "repo_push"
  method       = "webhook"
  title        = "import test"
  config       = { url = "https://httpbin.org/post" }
  event_config = {}
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_repository_notification.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				GenerateConfig:    true,
				ImportStateIdFunc: notificationImportID("quay_repository_notification.test"),
			},
		},
	})
}

// ---------------------------------------------------------------------------
// quay_repository_team_permission
// ---------------------------------------------------------------------------

func TestAccRepositoryTeamPermission_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-teamperm-cmd")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "import-perm-repo"
  visibility  = "private"
  description = ""
}
resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "import-team"
  role     = "member"
}
resource "quay_repository_team_permission" "test" {
  repository = quay_repository.test.id
  team_name  = quay_organization_team.test.name
  role       = "read"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_repository_team_permission.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRepositoryTeamPermission_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-teamperm-blk")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "import-perm-repo"
  visibility  = "private"
  description = ""
}
resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "import-team"
  role     = "member"
}
resource "quay_repository_team_permission" "test" {
  repository = quay_repository.test.id
  team_name  = quay_organization_team.test.name
  role       = "read"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_repository_team_permission.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRepositoryTeamPermission_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-teamperm-gen")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "import-perm-repo"
  visibility  = "private"
  description = ""
}
resource "quay_organization_team" "test" {
  org_name = quay_organization.test.name
  name     = "import-team"
  role     = "member"
}
resource "quay_repository_team_permission" "test" {
  repository = quay_repository.test.id
  team_name  = quay_organization_team.test.name
  role       = "read"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:    "quay_repository_team_permission.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
			},
		},
	})
}

// ---------------------------------------------------------------------------
// quay_repository_user_permission
// ---------------------------------------------------------------------------

func TestAccRepositoryUserPermission_importCommand(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-userperm-cmd")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "import-userperm-repo"
  visibility  = "private"
  description = ""
}
resource "quay_repository_user_permission" "test" {
  repository = quay_repository.test.id
  username   = "quayadmin"
  role       = "read"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_repository_user_permission.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRepositoryUserPermission_importBlock(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-userperm-blk")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "import-userperm-repo"
  visibility  = "private"
  description = ""
}
resource "quay_repository_user_permission" "test" {
  repository = quay_repository.test.id
  username   = "quayadmin"
  role       = "read"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:      "quay_repository_user_permission.test",
				ImportState:       true,
				ImportStateKind:   resource.ImportBlockWithID,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRepositoryUserPermission_generateConfig(t *testing.T) {
	skipIfNoAcc(t)
	orgName := rname("tf-acc-import-userperm-gen")
	cfg := providerConfig() + fmt.Sprintf(`
resource "quay_organization" "test" {
  name  = %q
  email = "import@example.com"
}
resource "quay_repository" "test" {
  namespace   = quay_organization.test.name
  name        = "import-userperm-repo"
  visibility  = "private"
  description = ""
}
resource "quay_repository_user_permission" "test" {
  repository = quay_repository.test.id
  username   = "quayadmin"
  role       = "read"
}
`, orgName)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				ResourceName:    "quay_repository_user_permission.test",
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,
				GenerateConfig:  true,
			},
		},
	})
}
