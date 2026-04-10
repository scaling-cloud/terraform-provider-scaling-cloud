terraform {
  required_providers {
    scaling = {
      source = "scaling-cloud/scaling-cloud"
    }
  }
}

# Configure the provider using the SCALING_CLOUD_API_KEY environment variable
provider "scaling" {}
