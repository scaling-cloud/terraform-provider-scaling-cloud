package oncall_schedule

import (
	"testing"

	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

func TestMapToModel(t *testing.T) {
	t.Parallel()

	got := mapToModel(&client.OncallSchedule{
		ID:        "sch_1",
		OrgID:     "org_1",
		Name:      "Primary",
		Timezone:  "Europe/London",
		CreatedAt: "t1",
		UpdatedAt: "t2",
	})

	if got.ID.ValueString() != "sch_1" || got.Name.ValueString() != "Primary" {
		t.Errorf("got %q/%q, want sch_1/Primary", got.ID.ValueString(), got.Name.ValueString())
	}
	if got.Timezone.ValueString() != "Europe/London" {
		t.Errorf("Timezone = %q, want Europe/London", got.Timezone.ValueString())
	}
	if !got.Description.IsNull() {
		t.Errorf("Description = %v, want null", got.Description)
	}
}
