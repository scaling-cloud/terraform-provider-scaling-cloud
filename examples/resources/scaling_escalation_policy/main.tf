resource "scaling_oncall_schedule" "primary" {
  name     = "Primary On-Call"
  timezone = "America/New_York"

  layer {
    name                 = "Tier 1"
    rotation_type        = "weekly"
    rotation_length_days = 7
    handoff_time         = "09:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_abc", "user_def"]
  }
}

resource "scaling_oncall_schedule" "secondary" {
  name     = "Secondary On-Call"
  timezone = "America/New_York"

  layer {
    name                 = "Backup"
    rotation_type        = "weekly"
    rotation_length_days = 7
    handoff_time         = "09:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_ghi", "user_jkl"]
  }
}

resource "scaling_escalation_policy" "default" {
  name        = "Default Escalation"
  description = "Escalate from primary to secondary on-call"

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.primary.id
    escalate_after_seconds = 300
  }

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.secondary.id
    escalate_after_seconds = 600
  }
}
