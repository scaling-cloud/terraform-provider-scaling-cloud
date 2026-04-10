package testutil

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/provider"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"scaling": providerserver.NewProtocol6WithError(provider.New("test")()),
}

func TestAccPreCheck(t *testing.T) {
	t.Helper()
	if os.Getenv("SCALING_CLOUD_API_KEY") == "" {
		t.Skip("SCALING_CLOUD_API_KEY not set, skipping acceptance test")
	}
}

func ProviderConfig() string {
	return `provider "scaling" {}`
}
