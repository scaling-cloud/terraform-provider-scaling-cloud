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

resource "scaling_escalation_policy" "payments" {
  name = "Payments Escalation"

  step {
    target_type            = "schedule"
    target_id              = scaling_oncall_schedule.primary.id
    escalate_after_seconds = 300
  }
}

# A routing policy carries exactly four rules — one per severity. The
# escalation_policy_id selects which escalation policy handles paging outcomes;
# omit it to fall back to the alert's component default.
resource "scaling_routing_policy" "payments" {
  name        = "Payments Routing"
  description = "Stricter routing for payment incidents"

  rule {
    severity             = "critical"
    outcome              = "incident"
    escalation_policy_id = scaling_escalation_policy.payments.id
  }

  rule {
    severity = "high"
    outcome  = "provisional_page"
  }

  rule {
    severity = "medium"
    outcome  = "notification"
  }

  rule {
    severity = "low"
    outcome  = "drop"
  }
}
