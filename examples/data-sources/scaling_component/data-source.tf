# Reference a component created out-of-band by name, for example to attach an
# inbound integration or report against it.
data "scaling_component" "checkout" {
  name = "checkout"
}

output "checkout_component_id" {
  value = data.scaling_component.checkout.id
}
