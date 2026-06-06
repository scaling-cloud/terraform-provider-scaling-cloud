# A component models a part of your system. Aliases are alternate names that
# inbound alerts can resolve this component by — the set is replaced wholesale
# on every apply.
resource "scaling_component" "payments_api" {
  name        = "Payments API"
  description = "Handles checkout and billing requests"

  aliases = [
    "payments",
    "checkout",
    "billing-service",
  ]
}
