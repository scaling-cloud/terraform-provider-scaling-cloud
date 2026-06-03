resource "scaling_oncall_schedule" "primary" {
  name     = "Primary On-Call"
  timezone = "America/New_York"

  layer {
    name                 = "Engineers"
    rotation_type        = "weekly"
    rotation_length_days = 7
    handoff_time         = "09:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_abc", "user_def"]
  }
}

resource "scaling_oncall_schedule" "management" {
  name     = "Management"
  timezone = "America/New_York"

  layer {
    name                 = "Managers"
    rotation_type        = "weekly"
    rotation_length_days = 7
    handoff_time         = "09:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_ghi", "user_jkl"]
  }
}

# Steps escalate in the order written. Each step pages its target, then waits
# escalate_after_seconds before escalating to the next step.
resource "scaling_escalation_policy" "default" {
  name        = "Default Escalation"
  description = "Page the primary rotation, then management after 5 minutes"

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.primary.id
    escalate_after_seconds = 300
  }

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.management.id
    escalate_after_seconds = 600
  }
}
