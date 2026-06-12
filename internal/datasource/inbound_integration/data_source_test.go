package inbound_integration

import (
	"testing"

	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

func TestMapToModel(t *testing.T) {
	t.Parallel()

	rp := "rp_1"
	got := mapToModel(&client.InboundIntegration{
		ID:              "int_1",
		OrgID:           "org_1",
		Name:            "Datadog",
		ComponentID:     "cmp_1",
		RoutingPolicyID: &rp,
		CreatedAt:       "2026-01-05T08:00:00.000Z",
		UpdatedAt:       "2026-03-20T16:45:00.000Z",
	})

	if got.ID.ValueString() != "int_1" {
		t.Errorf("ID = %q, want int_1", got.ID.ValueString())
	}
	if got.Name.ValueString() != "Datadog" {
		t.Errorf("Name = %q, want Datadog", got.Name.ValueString())
	}
	if got.ComponentID.ValueString() != "cmp_1" {
		t.Errorf("ComponentID = %q, want cmp_1", got.ComponentID.ValueString())
	}
	if got.RoutingPolicyID.ValueString() != "rp_1" {
		t.Errorf("RoutingPolicyID = %q, want rp_1", got.RoutingPolicyID.ValueString())
	}
}

func TestMapToModelNullRoutingPolicy(t *testing.T) {
	t.Parallel()

	// A null routing_policy_id (org default) must surface as a Terraform null,
	// not an empty string, so the attribute reads as unset.
	got := mapToModel(&client.InboundIntegration{ID: "int_1", RoutingPolicyID: nil})
	if !got.RoutingPolicyID.IsNull() {
		t.Errorf("RoutingPolicyID = %v, want null", got.RoutingPolicyID)
	}
}
