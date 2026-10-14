# Terraform Provider for Quay

A Terraform provider for managing resources on [Quay.io](https://quay.io) or a self-hosted [Red Hat Quay](https://www.redhat.com/en/technologies/cloud-computing/quay) instance.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Usage

```hcl
terraform {
  required_providers {
    quay = {
      source  = "registry.terraform.io/lannuttia/quay"
      version = "~> 0.1"
    }
  }
}

provider "quay" {
  # url defaults to https://quay.io
  # Can also be set via QUAY_URL environment variable
  url = "https://quay.io"

  # Can also be set via QUAY_TOKEN environment variable
  token = var.quay_token
}
```

## Supported Resources

| Resource | Description |
|---|---|
| `quay_organization` | Manages a Quay organization |
| `quay_organization_application` | Manages an OAuth application in an organization |
| `quay_organization_robot` | Manages a robot account in an organization |
| `quay_organization_team` | Manages a team in an organization |
| `quay_organization_team_member` | Manages membership in an organization team |
| `quay_organization_prototype_permission` | Manages default permission prototypes for an organization |
| `quay_repository` | Manages a Quay repository |
| `quay_repository_notification` | Manages a notification (webhook) on a repository |
| `quay_repository_team_permission` | Manages a team's permission on a repository |
| `quay_repository_user_permission` | Manages a user's permission on a repository |

## Supported Data Sources

| Data Source | Description |
|---|---|
| `data.quay_organization` | Reads information about a Quay organization |
| `data.quay_repository` | Reads information about a Quay repository |
| `data.quay_user` | Reads information about a Quay user |

## Development

### Building

```shell
make build
```

### Installing Locally

```shell
make install
```

This installs the provider binary into `~/.terraform.d/plugins/` so you can use it with a local Terraform configuration.

### Running Acceptance Tests

Acceptance tests require a real Quay instance and a valid API token:

```shell
export QUAY_URL=https://quay.io   # or your self-hosted URL
export QUAY_TOKEN=your_token_here
make testacc
```
