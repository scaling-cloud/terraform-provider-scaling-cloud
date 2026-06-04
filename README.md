# Terraform Provider for Scaling Cloud

The Scaling Cloud provider allows [Terraform](https://www.terraform.io) to manage resources on the [Scaling Cloud](https://scaling.cloud) incident management platform.

## Requirements

- [Terraform](https://www.terraform.io/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (for building from source)

## Installation

The provider is published on the [Terraform Registry](https://registry.terraform.io/providers/scaling-cloud/scaling-cloud/latest). Declare it in your configuration and run `terraform init` to install it automatically:

```hcl
terraform {
  required_providers {
    scaling = {
      source  = "scaling-cloud/scaling-cloud"
      version = "~> 0.1"
    }
  }
}
```

## Documentation

Full provider and resource documentation lives on the Registry:

- [Provider configuration](https://registry.terraform.io/providers/scaling-cloud/scaling-cloud/latest/docs)
- [`scaling_oncall_schedule`](https://registry.terraform.io/providers/scaling-cloud/scaling-cloud/latest/docs/resources/oncall_schedule)
- [`scaling_escalation_policy`](https://registry.terraform.io/providers/scaling-cloud/scaling-cloud/latest/docs/resources/escalation_policy)
- [`scaling_routing_policy`](https://registry.terraform.io/providers/scaling-cloud/scaling-cloud/latest/docs/resources/routing_policy)
- [`scaling_working_hours`](https://registry.terraform.io/providers/scaling-cloud/scaling-cloud/latest/docs/resources/working_hours)

The same content is generated into the [`docs/`](docs/) directory and rendered by the Registry.

## Usage

```hcl
terraform {
  required_providers {
    scaling = {
      source = "scaling-cloud/scaling-cloud"
    }
  }
}

provider "scaling" {
  # API key can also be set via the SCALING_CLOUD_API_KEY environment variable
  api_key = var.scaling_cloud_api_key
}

resource "scaling_oncall_schedule" "primary" {
  name     = "Primary On-Call"
  timezone = "America/New_York"

  layer {
    name                 = "Tier 1"
    rotation_type        = "weekly"
    rotation_length_days = 7
    handoff_time         = "09:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_1", "user_2"]
  }
}

resource "scaling_escalation_policy" "default" {
  name = "Default Escalation"

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.primary.id
    escalate_after_seconds = 300
  }
}
```

## Authentication

The provider authenticates using API keys. Configure via:

- Provider attribute: `api_key`
- Environment variable: `SCALING_CLOUD_API_KEY` (recommended)

```bash
export SCALING_CLOUD_API_KEY="your-api-key"
terraform plan
```

## Resources

- `scaling_oncall_schedule` - Manages on-call schedules with rotation layers
- `scaling_escalation_policy` - Manages escalation policies with ordered steps
- `scaling_routing_policy` - Manages routing policies mapping each alert severity to an outcome
- `scaling_working_hours` - Manages reusable Working Hours sets (timezone + weekly windows) for follow-the-sun escalation step conditions

## Building the Provider

```bash
go build -o terraform-provider-scaling-cloud .
```

## Development

### Running Tests

Unit tests:

```bash
go test -v -race ./...
```

Acceptance tests (requires API access):

```bash
export TF_ACC=1
export SCALING_CLOUD_API_KEY="your-test-api-key"
go test -v -count=1 -timeout 10m ./internal/...
```

### Linting

```bash
golangci-lint run
```

### Generating Documentation

The `docs/` directory is generated from the provider schema and the `examples/`
directory using [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs),
which is pinned as a Go tool dependency. Regenerate and commit it whenever a
resource schema or example changes:

```bash
go generate ./...
```

CI fails if `docs/` is out of date.

## Releasing

Releases are cut by pushing a semver tag. The `release` workflow cross-compiles
the binaries with [GoReleaser](https://goreleaser.com), produces a `SHA256SUMS`
file, and signs it with the GPG key stored in the `GPG_PRIVATE_KEY` and
`PASSPHRASE` repository secrets. The Terraform Registry ingests the resulting
GitHub release automatically.

```bash
git tag v0.1.0
git push origin v0.1.0
```

Publishing to the Registry requires a one-time setup: connect the GitHub
repository at [registry.terraform.io](https://registry.terraform.io/publish) and
register the public half of the signing key under the namespace's GPG keys.

## License

[MPL-2.0](LICENSE)
