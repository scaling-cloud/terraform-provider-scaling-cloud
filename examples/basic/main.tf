terraform {
  required_providers {
    scaling = {
      source = "scaling-cloud/scaling-cloud"
    }
  }
}

# Configure the provider using the SCALING_CLOUD_API_KEY environment variable
provider "scaling" {}

resource "scaling_oncall_schedule" "example" {
  name     = "Engineering On-Call"
  timezone = "America/New_York"
}

resource "scaling_escalation_policy" "example" {
  name = "Default Escalation"

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.example.id
    escalate_after_seconds = 300
  }
}
