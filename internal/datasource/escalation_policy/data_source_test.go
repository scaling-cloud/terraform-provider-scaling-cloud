package escalation_policy

import (
	"testing"

	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

func TestMapToModel(t *testing.T) {
	t.Parallel()

	desc := "critical path"
	got := mapToModel(&client.EscalationPolicy{
		ID:          "ep_1",
		OrgID:       "org_1",
		Name:        "Critical",
		Description: &desc,
		CreatedAt:   "t1",
		UpdatedAt:   "t2",
	})

	if got.ID.ValueString() != "ep_1" || got.Name.ValueString() != "Critical" {
		t.Errorf("got %q/%q, want ep_1/Critical", got.ID.ValueString(), got.Name.ValueString())
	}
	if got.Description.ValueString() != "critical path" {
		t.Errorf("Description = %q, want critical path", got.Description.ValueString())
	}
}

func TestMapToModelNullDescription(t *testing.T) {
	t.Parallel()

	got := mapToModel(&client.EscalationPolicy{ID: "ep_1", Description: nil})
	if !got.Description.IsNull() {
		t.Errorf("Description = %v, want null", got.Description)
	}
}
