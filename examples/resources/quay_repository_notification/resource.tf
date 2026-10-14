# Webhook triggered on every push
resource "quay_repository_notification" "push_webhook" {
  repository = "my-org/my-app"
  event      = "repo_push"
  method     = "webhook"
  title      = "Push notification"

  config = {
    url = "https://hooks.example.com/quay"
  }

  event_config = {}
}

# Slack notification for build failures
resource "quay_repository_notification" "build_failure_slack" {
  repository = "my-org/my-app"
  event      = "build_failure"
  method     = "slack"
  title      = "Build failure alert"

  config = {
    url = "https://hooks.slack.com/services/T00000/B00000/XXXXXXXXXXXX"
  }

  event_config = {}
}

# Vulnerability notification for High and above
resource "quay_repository_notification" "vuln_webhook" {
  repository = "my-org/my-app"
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
