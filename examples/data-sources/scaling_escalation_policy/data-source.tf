# Reference an escalation policy created out-of-band by name, for use as a
# routing policy rule's escalation_policy_id.
data "scaling_escalation_policy" "critical" {
  name = "Critical"
}

output "critical_escalation_policy_id" {
  value = data.scaling_escalation_policy.critical.id
}
