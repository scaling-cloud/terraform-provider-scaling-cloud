terraform {
  required_providers {
    scaling = {
      source = "scaling-cloud/scaling-cloud"
    }
  }
}

provider "scaling" {}

# On-call schedules for different teams

resource "scaling_oncall_schedule" "platform" {
  name        = "Platform Team"
  description = "Platform engineering on-call rotation"
  timezone    = "America/New_York"

  layer {
    name                 = "Primary"
    rotation_type        = "weekly"
    rotation_length_days = 7
    handoff_time         = "09:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_alice", "user_bob", "user_carol"]
  }

  layer {
    name                 = "Weekend Override"
    rotation_type        = "custom"
    rotation_length_days = 2
    handoff_time         = "18:00"
    effective_from       = "2025-01-04T00:00:00Z"
    participant_ids      = ["user_dave", "user_eve"]
  }
}

resource "scaling_oncall_schedule" "backend" {
  name        = "Backend Team"
  description = "Backend services on-call rotation"
  timezone    = "Europe/London"

  layer {
    name                 = "Primary"
    rotation_type        = "daily"
    rotation_length_days = 1
    handoff_time         = "08:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_frank", "user_grace"]
  }
}

resource "scaling_oncall_schedule" "management" {
  name        = "Engineering Management"
  description = "Management escalation for critical incidents"
  timezone    = "America/New_York"

  layer {
    name                 = "Managers"
    rotation_type        = "weekly"
    rotation_length_days = 7
    handoff_time         = "09:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_heidi", "user_ivan"]
  }
}

# Escalation policies

resource "scaling_escalation_policy" "platform_incidents" {
  name        = "Platform Incidents"
  description = "Platform team escalation: primary -> backend -> management"

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.platform.id
    escalate_after_seconds = 300
  }

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.backend.id
    escalate_after_seconds = 600
  }

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.management.id
    escalate_after_seconds = 900
  }
}

resource "scaling_escalation_policy" "backend_incidents" {
  name        = "Backend Incidents"
  description = "Backend team escalation: backend -> management"

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.backend.id
    escalate_after_seconds = 300
  }

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.management.id
    escalate_after_seconds = 600
  }
}
