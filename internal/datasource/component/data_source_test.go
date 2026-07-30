package component

import (
	"context"
	"slices"
	"testing"

	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/datasource/lookup"
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

func TestFilterByAlias_matchesComponent(t *testing.T) {
	t.Parallel()

	components := []client.Component{
		{ID: "cmp_1", Name: "billing", Aliases: []string{"payments", "checkout"}},
		{ID: "cmp_2", Name: "auth", Aliases: []string{"login"}},
		{ID: "cmp_3", Name: "other", Aliases: nil},
	}

	got := filterByAlias(components, "payments")
	if len(got) != 1 || got[0].ID != "cmp_1" {
		t.Fatalf("filterByAlias(payments) = %d components, want 1 (cmp_1)", len(got))
	}
}

func TestFilterByAlias_noMatch(t *testing.T) {
	t.Parallel()

	components := []client.Component{
		{ID: "cmp_1", Name: "billing", Aliases: []string{"payments"}},
	}

	got := filterByAlias(components, "nonexistent")
	if len(got) != 0 {
		t.Fatalf("filterByAlias(nonexistent) = %d components, want 0", len(got))
	}
}

func TestFilterByAlias_multipleMatch(t *testing.T) {
	t.Parallel()

	components := []client.Component{
		{ID: "cmp_1", Name: "billing", Aliases: []string{"shared"}},
		{ID: "cmp_2", Name: "auth", Aliases: []string{"shared"}},
	}

	got := filterByAlias(components, "shared")
	if len(got) != 2 {
		t.Fatalf("filterByAlias(shared) = %d components, want 2", len(got))
	}
	ids := []string{got[0].ID, got[1].ID}
	if !slices.Contains(ids, "cmp_1") || !slices.Contains(ids, "cmp_2") {
		t.Errorf("missing expected component IDs in %v", ids)
	}
}

func TestFilterByAlias_emptyInput(t *testing.T) {
	t.Parallel()

	got := filterByAlias(nil, "anything")
	if len(got) != 0 {
		t.Fatalf("filterByAlias(nil) = %d components, want 0", len(got))
	}
}

func TestFilterByAlias_nilAliasesList(t *testing.T) {
	t.Parallel()

	components := []client.Component{
		{ID: "cmp_1", Name: "billing", Aliases: nil},
	}

	got := filterByAlias(components, "anything")
	if len(got) != 0 {
		t.Fatalf("filterByAlias on nil aliases = %d components, want 0", len(got))
	}
}

func TestMapByAliasAndName(t *testing.T) {
	t.Parallel()

	// Simulates the full Read path: filter by alias then lookup by name.
	components := []client.Component{
		{ID: "cmp_1", Name: "billing", Aliases: []string{"payments"}},
		{ID: "cmp_2", Name: "billing", Aliases: []string{"invoices"}},  // same name, different alias
		{ID: "cmp_3", Name: "auth", Aliases: []string{"payments"}},
	}

	filtered := filterByAlias(components, "payments")
	if len(filtered) != 2 {
		t.Fatalf("filterByAlias(payments) = %d, want 2", len(filtered))
	}

	match, err := lookup.One(filtered, "component", "name", "billing", func(c client.Component) string { return c.Name })
	if err != nil {
		t.Fatalf("lookup.One(billing) after alias filter: %v", err)
	}
	if match.ID != "cmp_1" {
		t.Errorf("expected cmp_1, got %s", match.ID)
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
