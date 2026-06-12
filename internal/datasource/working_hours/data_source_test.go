package working_hours

import (
	"testing"

	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

func TestMapToModel(t *testing.T) {
	t.Parallel()

	got := mapToModel(&client.WorkingHours{
		ID:        "wh_1",
		OrgID:     "org_1",
		Name:      "Business Hours",
		Timezone:  "Europe/London",
		CreatedAt: "t1",
		UpdatedAt: "t2",
	})

	if got.ID.ValueString() != "wh_1" || got.Name.ValueString() != "Business Hours" {
		t.Errorf("got %q/%q, want wh_1/Business Hours", got.ID.ValueString(), got.Name.ValueString())
	}
	if got.Timezone.ValueString() != "Europe/London" {
		t.Errorf("Timezone = %q, want Europe/London", got.Timezone.ValueString())
	}
}
