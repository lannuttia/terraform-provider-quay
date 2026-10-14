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
  method     = "webhook"
  title      = "Push notification"

  config = {
    url = "https://hooks.example.com/quay"
  }

  event_config = {}
}

# Slack alert for build failures
resource "quay_repository_notification" "build_failure_slack" {
  repository = quay_repository.example.id
  event      = "build_failure"
  method     = "slack"
  title      = "Build failure alert"

  config = {
    url = "https://hooks.slack.com/services/T00000/B00000/XXXXXXXXXXXX"
  }

  event_config = {}
}

# Webhook for High severity (and above) vulnerability findings
resource "quay_repository_notification" "vuln_webhook" {
  repository = quay_repository.example.id
  event      = "vulnerability_found"
  method     = "webhook"
  title      = "High+ vulnerability alert"

  config = {
    url = "https://hooks.example.com/vulns"
  }

  event_config = {
    level = "High"
  }
}
