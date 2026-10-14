resource "quay_organization_team" "developers" {
  org_name    = "my-org"
  name        = "developers"
  role        = "member"
  description = "Development team"
}

resource "quay_organization_team_member" "alice" {
  org_name    = "my-org"
  team_name   = quay_organization_team.developers.name
  member_name = "alice"
}

resource "quay_repository_team_permission" "developers_write" {
  repository = "my-org/my-app"
  team_name  = quay_organization_team.developers.name
  role       = "write"
}
