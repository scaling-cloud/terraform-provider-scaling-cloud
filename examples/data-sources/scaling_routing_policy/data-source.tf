# Reference a routing policy created out-of-band by name, for use in an inbound
# integration selector's routing_policy_id.
data "scaling_routing_policy" "default" {
  name = "Default"
}

output "default_routing_policy_id" {
  value = data.scaling_routing_policy.default.id
}
