resource "quay_organization" "example" {
  name  = "my-org"
  email = "admin@example.com"
}

resource "quay_repository" "example" {
  namespace  = quay_organization.example.name
  name       = "my-app"
  visibility = "private"
}

# Webhook triggered on every image push
resource "quay_repository_notification" "push_webhook" {
  repository = quay_repository.example.id
  event      = "repo_push"
  title      = "Push notification"

  webhook = {
    url = "https://hooks.example.com/quay"
  }
}

# Slack alert for build failures
resource "quay_repository_notification" "build_failure_slack" {
  repository = quay_repository.example.id
  event      = "build_failure"
  title      = "Build failure alert"

  slack = {
    url = "https://hooks.slack.com/services/T00000/B00000/XXXXXXXXXXXX"
  }
}

# Webhook for High severity (and above) vulnerability findings
resource "quay_repository_notification" "vuln_webhook" {
  repository = quay_repository.example.id
  event      = "vulnerability_found"
  title      = "High+ vulnerability alert"

  webhook = {
    url = "https://hooks.example.com/vulns"
  }

  vulnerability_config = {
    level = "High"
  }
}

# Email notification for main-branch build successes only
resource "quay_repository_notification" "main_build_email" {
  repository = quay_repository.example.id
  event      = "build_success"
  title      = "Main branch build succeeded"

  email = {
    email = "team@example.com"
  }

  build_config = {
    ref_regex = "refs/heads/main"
  }
}
