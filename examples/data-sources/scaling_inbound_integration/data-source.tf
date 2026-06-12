# Look up an inbound integration that was installed out-of-band (for example by
# connecting a source in the UI) so its selectors can be managed in Terraform.
data "scaling_inbound_integration" "datadog" {
  name = "Datadog"
}

resource "scaling_inbound_integration" "datadog" {
  integration_id = data.scaling_inbound_integration.datadog.id

  selector {
    routing_policy_id = scaling_routing_policy.critical.id

    matcher {
      key   = "env"
      value = "prod"
    }
  }
}
