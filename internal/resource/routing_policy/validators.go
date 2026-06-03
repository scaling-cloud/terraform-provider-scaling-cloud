package routing_policy

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ validator.List = &oneRulePerSeverityValidator{}

// oneRulePerSeverityValidator enforces that a routing policy's rule block
// contains exactly one rule per supported severity, so the four-rule set is
// complete and non-duplicated before the API is ever called.
type oneRulePerSeverityValidator struct{}

func (v *oneRulePerSeverityValidator) Description(_ context.Context) string {
	return "must contain exactly one rule per severity (critical, high, medium, low)"
}

func (v *oneRulePerSeverityValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v *oneRulePerSeverityValidator) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	elements := req.ConfigValue.Elements()
	seen := make(map[string]int, len(validSeverities))

	for _, element := range elements {
		obj, ok := element.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			continue
		}
		severityAttr, ok := obj.Attributes()["severity"]
		if !ok {
			continue
		}
		severity, ok := severityAttr.(types.String)
		if !ok || severity.IsNull() || severity.IsUnknown() {
			continue
		}
		seen[severity.ValueString()]++
	}

	var missing, duplicated []string
	for _, severity := range validSeverities {
		switch count := seen[severity]; {
		case count == 0:
			missing = append(missing, severity)
		case count > 1:
			duplicated = append(duplicated, severity)
		}
	}

	if len(missing) == 0 && len(duplicated) == 0 {
		return
	}

	var details []string
	if len(missing) > 0 {
		details = append(details, fmt.Sprintf("missing rules for: %s", strings.Join(missing, ", ")))
	}
	if len(duplicated) > 0 {
		details = append(details, fmt.Sprintf("duplicate rules for: %s", strings.Join(duplicated, ", ")))
	}

	resp.Diagnostics.AddAttributeError(
		req.Path,
		"Invalid routing rule set",
		fmt.Sprintf("A routing policy must define exactly one rule per severity (%s); %s.",
			strings.Join(validSeverities, ", "),
			strings.Join(details, "; "),
		),
	)
}
