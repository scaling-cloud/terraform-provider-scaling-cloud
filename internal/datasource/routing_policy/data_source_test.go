package routing_policy

import (
	"testing"

	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

func TestMapToModel(t *testing.T) {
	t.Parallel()

	got := mapToModel(&client.RoutingPolicy{
		ID:        "rp_1",
		OrgID:     "org_1",
		Name:      "Default",
		IsDefault: true,
		CreatedAt: "t1",
		UpdatedAt: "t2",
	})

	if got.ID.ValueString() != "rp_1" || got.Name.ValueString() != "Default" {
		t.Errorf("got %q/%q, want rp_1/Default", got.ID.ValueString(), got.Name.ValueString())
	}
	if !got.IsDefault.ValueBool() {
		t.Errorf("IsDefault = false, want true")
	}
	// A null description must surface as a Terraform null, not "".
	if !got.Description.IsNull() {
		t.Errorf("Description = %v, want null", got.Description)
	}
}
