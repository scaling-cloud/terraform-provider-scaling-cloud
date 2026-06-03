package routing_policy

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ruleObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"severity":             types.StringType,
			"outcome":              types.StringType,
			"escalation_policy_id": types.StringType,
		},
	}
}

func ruleObject(t *testing.T, severity string) types.Object {
	t.Helper()
	obj, diags := types.ObjectValue(ruleObjectType().AttrTypes, map[string]attr.Value{
		"severity":             types.StringValue(severity),
		"outcome":              types.StringValue("notification"),
		"escalation_policy_id": types.StringNull(),
	})
	if diags.HasError() {
		t.Fatalf("building rule object: %v", diags)
	}
	return obj
}

func runValidator(t *testing.T, severities []string) validator.ListResponse {
	t.Helper()
	elements := make([]attr.Value, len(severities))
	for i, s := range severities {
		elements[i] = ruleObject(t, s)
	}
	list, diags := types.ListValue(ruleObjectType(), elements)
	if diags.HasError() {
		t.Fatalf("building list: %v", diags)
	}

	req := validator.ListRequest{
		Path:        path.Root("rule"),
		ConfigValue: list,
	}
	resp := validator.ListResponse{}
	v := &oneRulePerSeverityValidator{}
	v.ValidateList(context.Background(), req, &resp)
	return resp
}

func TestOneRulePerSeverity_complete(t *testing.T) {
	t.Parallel()
	resp := runValidator(t, []string{"critical", "high", "medium", "low"})
	if resp.Diagnostics.HasError() {
		t.Fatalf("expected no error for complete set, got: %v", resp.Diagnostics)
	}
}

func TestOneRulePerSeverity_missing(t *testing.T) {
	t.Parallel()
	resp := runValidator(t, []string{"critical", "high", "medium"})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for missing severity")
	}
}

func TestOneRulePerSeverity_duplicate(t *testing.T) {
	t.Parallel()
	resp := runValidator(t, []string{"critical", "critical", "high", "medium"})
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error for duplicate severity")
	}
}
