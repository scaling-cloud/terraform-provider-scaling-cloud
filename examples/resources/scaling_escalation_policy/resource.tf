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

# Working Hours sets let steps page follow-the-sun. The condition is evaluated
# at the instant escalation reaches the step; an ineligible step is skipped
# (never held). If every step's condition fails, the last step pages anyway.
resource "scaling_working_hours" "uk_office" {
  name     = "UK office hours"
  timezone = "Europe/London"

  window {
    days  = [1, 2, 3, 4, 5]
    start = "09:00"
    end   = "17:00"
  }
}

# Steps escalate in the order written. Each step pages its target, then waits
# escalate_after_seconds before escalating to the next step. An optional
# condition block gates the step to a Working Hours set.
resource "scaling_escalation_policy" "default" {
  name        = "Default Escalation"
  description = "Page the primary rotation during UK hours, then management"

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.primary.id
    escalate_after_seconds = 300

    condition {
      working_hours_id = scaling_working_hours.uk_office.id
      when             = "during"
    }
  }

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.management.id
    escalate_after_seconds = 600
  }
}
