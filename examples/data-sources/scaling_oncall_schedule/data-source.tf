# Reference an on-call schedule created out-of-band by name, for use as an
# escalation step's target_id.
data "scaling_oncall_schedule" "primary" {
  name = "Primary On-Call"
}

output "primary_schedule_id" {
  value = data.scaling_oncall_schedule.primary.id
}
