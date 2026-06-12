package component

import (
	"context"
	"testing"

	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

func TestMapToModel(t *testing.T) {
	t.Parallel()

	desc := "checkout service"
	got, diags := mapToModel(context.Background(), &client.Component{
		ID:                "cmp_1",
		OrgID:             "org_1",
		Name:              "billing",
		Description:       &desc,
		Aliases:           []string{"payments", "checkout"},
		OperationalStatus: "operational",
		CreatedAt:         "t1",
		UpdatedAt:         "t2",
	})
	if diags.HasError() {
		t.Fatalf("mapToModel diags: %v", diags)
	}

	if got.ID.ValueString() != "cmp_1" || got.Name.ValueString() != "billing" {
		t.Errorf("got %q/%q, want cmp_1/billing", got.ID.ValueString(), got.Name.ValueString())
	}
	if got.OperationalStatus.ValueString() != "operational" {
		t.Errorf("OperationalStatus = %q, want operational", got.OperationalStatus.ValueString())
	}
	if len(got.Aliases.Elements()) != 2 {
		t.Errorf("Aliases len = %d, want 2", len(got.Aliases.Elements()))
	}
}

func TestMapToModelNilAliasesIsKnownEmptySet(t *testing.T) {
	t.Parallel()

	got, diags := mapToModel(context.Background(), &client.Component{ID: "cmp_1", Aliases: nil})
	if diags.HasError() {
		t.Fatalf("mapToModel diags: %v", diags)
	}
	// A component with no aliases must surface as a known empty set, not null.
	if got.Aliases.IsNull() {
		t.Errorf("Aliases IsNull, want a known empty set")
	}
	if len(got.Aliases.Elements()) != 0 {
		t.Errorf("Aliases len = %d, want 0", len(got.Aliases.Elements()))
	}
}
