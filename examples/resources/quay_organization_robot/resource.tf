resource "quay_organization_robot" "ci" {
  org_name    = "my-org"
  short_name  = "ci-bot"
  description = "Robot account used by CI/CD pipelines"
}

# Grant the robot write access to a repository
resource "quay_repository_user_permission" "ci_write" {
  repository = "my-org/my-app"
  username   = quay_organization_robot.ci.full_name
  role       = "write"
}
