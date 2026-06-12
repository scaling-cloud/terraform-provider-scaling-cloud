# Reference a working-hours set created out-of-band by name, for use in an
# escalation step condition's working_hours_id.
data "scaling_working_hours" "business" {
  name = "Business Hours"
}

output "business_hours_id" {
  value = data.scaling_working_hours.business.id
}
