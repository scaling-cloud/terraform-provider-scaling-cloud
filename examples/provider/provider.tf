terraform {
  required_providers {
    scaling = {
      source = "scaling-cloud/scaling-cloud"
    }
  }
}

# The API key can also be supplied via the SCALING_CLOUD_API_KEY environment
# variable (recommended), in which case the api_key argument can be omitted.
provider "scaling" {
  api_key = var.scaling_cloud_api_key
}
