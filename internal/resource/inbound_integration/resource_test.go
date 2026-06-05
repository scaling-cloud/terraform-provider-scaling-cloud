package inbound_integration

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

func TestSelectorsToInputPreservesOrderAndMatchers(t *testing.T) {
	t.Parallel()

	models := []SelectorModel{
		{
			RoutingPolicyID: types.StringValue("rp_critical"),
			Matchers: []MatcherModel{
				{Key: types.StringValue("service"), Value: types.StringValue("payments")},
				{Key: types.StringValue("env"), Value: types.StringValue("prod")},
			},
		},
		{
			RoutingPolicyID: types.StringValue("rp_low"),
			Matchers: []MatcherModel{
				{Key: types.StringValue("env"), Value: types.StringValue("staging")},
			},
		},
	}

	got := selectorsToInput(models)

	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	// Ordering is meaningful and must be preserved exactly as authored.
	if got[0].RoutingPolicyID != "rp_critical" || got[1].RoutingPolicyID != "rp_low" {
		t.Errorf("order = [%q %q], want [rp_critical rp_low]", got[0].RoutingPolicyID, got[1].RoutingPolicyID)
	}
	if len(got[0].Matchers) != 2 {
		t.Fatalf("got[0].Matchers len = %d, want 2", len(got[0].Matchers))
	}
	if got[0].Matchers[0].Key != "service" || got[0].Matchers[0].Value != "payments" {
		t.Errorf("got[0].Matchers[0] = %+v, want service/payments", got[0].Matchers[0])
	}
	if got[0].Matchers[1].Key != "env" {
		t.Errorf("got[0].Matchers[1].Key = %q, want env", got[0].Matchers[1].Key)
	}
}

func TestSelectorsToInputNonNilWhenEmpty(t *testing.T) {
	t.Parallel()

	got := selectorsToInput(nil)
	if got == nil {
		t.Errorf("got nil, want a non-nil empty slice so the PUT clears selectors")
	}
	if len(got) != 0 {
		t.Errorf("len = %d, want 0", len(got))
	}
}

func TestSelectorsToStateRoundTrip(t *testing.T) {
	t.Parallel()

	apiSelectors := []client.RoutingSelector{
		{
			RoutingPolicyID: "rp_critical",
			Matchers: []client.SelectorMatcher{
				{Key: "service", Value: "payments"},
				{Key: "env", Value: "prod"},
			},
		},
		{
			RoutingPolicyID: "rp_low",
			Matchers:        []client.SelectorMatcher{{Key: "env", Value: "staging"}},
		},
	}

	got := selectorsToState(apiSelectors)

	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].RoutingPolicyID.ValueString() != "rp_critical" {
		t.Errorf("got[0].RoutingPolicyID = %q, want rp_critical", got[0].RoutingPolicyID.ValueString())
	}
	if len(got[0].Matchers) != 2 || got[0].Matchers[0].Key.ValueString() != "service" {
		t.Errorf("got[0].Matchers = %+v, unexpected", got[0].Matchers)
	}
	// Round-trip back through input keeps ordering and content stable.
	back := selectorsToInput(got)
	if back[1].RoutingPolicyID != "rp_low" {
		t.Errorf("round-trip order broken: %q", back[1].RoutingPolicyID)
	}
}
