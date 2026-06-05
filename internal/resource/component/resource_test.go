package component

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

func TestAliasesToInputSortedAndNonNil(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	set, diags := types.SetValueFrom(ctx, types.StringType, []string{"checkout", "payments", "billing"})
	if diags.HasError() {
		t.Fatalf("building set: %v", diags)
	}

	got, diags := aliasesToInput(ctx, set)
	if diags.HasError() {
		t.Fatalf("aliasesToInput: %v", diags)
	}

	// A set has no inherent order; aliasesToInput sorts so the request body is
	// deterministic and PUT-replacement diffs stay stable.
	want := []string{"billing", "checkout", "payments"}
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d (%v)", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestAliasesToInputNullBecomesEmptySlice(t *testing.T) {
	t.Parallel()

	got, diags := aliasesToInput(context.Background(), types.SetNull(types.StringType))
	if diags.HasError() {
		t.Fatalf("aliasesToInput: %v", diags)
	}
	// Full-replacement: a null/absent aliases must serialize as an empty slice,
	// never nil, so the server clears any prior aliases.
	if got == nil {
		t.Errorf("got nil, want a non-nil empty slice")
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

func TestAliasesToStateRoundTrip(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	comp := &client.Component{
		ID:      "cmp_1",
		Aliases: []string{"payments", "checkout"},
	}

	set, diags := aliasesToState(ctx, comp.Aliases)
	if diags.HasError() {
		t.Fatalf("aliasesToState: %v", diags)
	}

	var out []string
	if d := set.ElementsAs(ctx, &out, false); d.HasError() {
		t.Fatalf("ElementsAs: %v", d)
	}
	if len(out) != 2 {
		t.Fatalf("len = %d, want 2 (%v)", len(out), out)
	}
}

func TestAliasesToStateEmptyIsKnownEmptySet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	set, diags := aliasesToState(ctx, []string{})
	if diags.HasError() {
		t.Fatalf("aliasesToState: %v", diags)
	}
	if set.IsNull() {
		t.Errorf("set IsNull, want a known empty set")
	}
	if len(set.Elements()) != 0 {
		t.Errorf("len = %d, want 0", len(set.Elements()))
	}
}
