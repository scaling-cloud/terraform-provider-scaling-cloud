# Terraform Provider for Scaling Cloud

The Scaling Cloud provider allows [Terraform](https://www.terraform.io) to manage resources on the [Scaling Cloud](https://scaling.cloud) incident management platform.

## Requirements

- [Terraform](https://www.terraform.io/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (for building from source)

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

## License

[MPL-2.0](LICENSE)
