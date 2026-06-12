# Resolve a user by email to their id, for use in fields that reference users
# such as an on-call schedule layer's participant_ids.
data "scaling_user" "lead" {
  email = "lead@example.com"
}

resource "scaling_oncall_schedule" "primary" {
  name     = "Primary On-Call"
  timezone = "Europe/London"
}

# data.scaling_user.lead.id can now be passed wherever a user id is required.
output "lead_user_id" {
  value = data.scaling_user.lead.id
}
