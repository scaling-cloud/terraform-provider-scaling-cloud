# Reference a component created out-of-band by name, for example to attach an
# inbound integration or report against it.
data "scaling_component" "checkout" {
  name = "checkout"
}

output "checkout_component_id" {
  value = data.scaling_component.checkout.id
}

# Optionally narrow the lookup by alias to disambiguate when multiple
# components share the same name but have different aliases.
data "scaling_component" "checkout_by_alias" {
  name  = "checkout"
  alias = "payments"
}

output "checkout_by_alias_component_id" {
  value = data.scaling_component.checkout_by_alias.id
}
