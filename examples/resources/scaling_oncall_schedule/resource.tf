resource "scaling_oncall_schedule" "primary" {
  name        = "Primary On-Call"
  description = "Primary on-call rotation for the platform team"
  timezone    = "America/New_York"

  layer {
    name                 = "Tier 1"
    rotation_type        = "weekly"
    rotation_length_days = 7
    handoff_time         = "09:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_abc", "user_def"]
  }

  layer {
    name                 = "Tier 2 - Backup"
    rotation_type        = "weekly"
    rotation_length_days = 7
    handoff_time         = "09:00"
    effective_from       = "2025-01-01T00:00:00Z"
    participant_ids      = ["user_ghi", "user_jkl"]
  }
}
