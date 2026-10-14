---
page_title: "Quay Provider"
description: |-
  Use the Quay provider to manage resources on Quay.io or a self-hosted
  Red Hat Quay instance.
---

# Quay Provider

The Quay provider manages resources on [Quay.io](https://quay.io) or a
self-hosted [Red Hat Quay](https://www.redhat.com/en/technologies/cloud-computing/quay)
instance via the Quay REST API.

## Authentication

The provider requires an OAuth bearer token. Generate one in the Quay UI under
**Account Settings → User Settings → Docker CLI Password** (or via an
OAuth application if managing org-owned resources).

The token can be supplied in the provider configuration or via the
`QUAY_TOKEN` environment variable.

## Example Usage

```terraform
terraform {
  required_providers {
    quay = {
      source  = "registry.terraform.io/lannuttia/quay"
      version = "~> 0.1"
    }
  }
}

provider "quay" {
  # Defaults to https://quay.io. Set to your self-hosted instance URL if needed.
  # Can also be set via the QUAY_URL environment variable.
  url = "https://quay.io"

  # Required. Can also be set via the QUAY_TOKEN environment variable.
  token = var.quay_token
}
```

## Argument Reference

- `url` - (Optional) The base URL of the Quay instance. Defaults to
  `https://quay.io`. Can also be set via `QUAY_URL`.
- `token` - (Required) OAuth bearer token. Can also be set via `QUAY_TOKEN`.
