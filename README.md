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

## Publishing to the Terraform Registry

### One-time setup

#### 1. Generate a GPG signing key (RSA — the Registry rejects ECC keys)

```shell
gpg --full-generate-key
# Choose: (1) RSA and RSA, 4096 bits
```

#### 2. Add the public key to the Terraform Registry

```shell
# Copy this output
gpg --armor --export your-email@example.com
```

Go to **https://registry.terraform.io/settings/gpg-keys** and paste it.

#### 3. Add secrets to your GitHub repository

Go to **Settings → Secrets and variables → Actions** and add:

| Secret | Value |
|---|---|
| `GPG_PRIVATE_KEY` | `gpg --armor --export-secret-keys YOUR_KEY_ID` |
| `PASSPHRASE` | The passphrase you set when generating the key |

#### 4. Register the provider on the Terraform Registry

1. Sign in to [registry.terraform.io](https://registry.terraform.io) with GitHub.
2. Click **Publish → Provider** in the top-right navigation.
3. Select your GitHub namespace and the `terraform-provider-quay` repository.
4. Accept the Terms of Use and click **Publish Provider**.

The Registry will create a webhook on your repo. Every future GitHub Release
is automatically picked up.

### Publishing a release

```shell
git tag v0.1.0
git push origin v0.1.0
```

The `release.yml` workflow fires, GoReleaser builds binaries for all
platforms, signs the checksums with your GPG key, and creates a GitHub
Release. The Registry webhook ingests it within a few minutes.

---

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
