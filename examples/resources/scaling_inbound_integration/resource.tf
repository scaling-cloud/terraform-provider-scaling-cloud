resource "scaling_routing_policy" "critical" {
  name = "Critical Routing"

  rule {
    severity = "critical"
    outcome  = "incident"
  }
  rule {
    severity = "high"
    outcome  = "incident"
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

resource "scaling_routing_policy" "low_noise" {
  name = "Low Noise Routing"

  rule {
    severity = "critical"
    outcome  = "notification"
  }
  rule {
    severity = "high"
    outcome  = "notification"
  }
  rule {
    severity = "medium"
    outcome  = "drop"
  }
  rule {
    severity = "low"
    outcome  = "drop"
  }
}

# The inbound integration itself is provisioned out-of-band (for example by
# connecting a source in the UI). This resource owns only its ordered Routing
# Selectors: the first row whose matchers all match an alert wins, so order is
# significant. A miss falls through to the integration's default policy.
resource "scaling_inbound_integration" "datadog" {
  integration_id = "int_01h9xyz"

  selector {
    matcher {
      key   = "service"
      value = "payments"
    }
    matcher {
      key   = "env"
      value = "prod"
    }
    routing_policy_id = scaling_routing_policy.critical.id
  }

  selector {
    matcher {
      key   = "env"
      value = "staging"
    }
    routing_policy_id = scaling_routing_policy.low_noise.id
  }
}
