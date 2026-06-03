package routing_policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

var (
	_ resource.Resource              = &RoutingPolicyResource{}
	_ resource.ResourceWithConfigure = &RoutingPolicyResource{}
)

var validSeverities = []string{"critical", "high", "medium", "low"}

type RoutingPolicyResource struct {
	client *client.ScalingClient
}

type RoutingPolicyModel struct {
	ID          types.String       `tfsdk:"id"`
	OrgID       types.String       `tfsdk:"org_id"`
	Name        types.String       `tfsdk:"name"`
	Description types.String       `tfsdk:"description"`
	IsDefault   types.Bool         `tfsdk:"is_default"`
	CreatedAt   types.String       `tfsdk:"created_at"`
	UpdatedAt   types.String       `tfsdk:"updated_at"`
	Rules       []RoutingRuleModel `tfsdk:"rule"`
}

type RoutingRuleModel struct {
	Severity           types.String `tfsdk:"severity"`
	Outcome            types.String `tfsdk:"outcome"`
	EscalationPolicyID types.String `tfsdk:"escalation_policy_id"`
}

func NewResource() resource.Resource {
	return &RoutingPolicyResource{}
}

func (r *RoutingPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_policy"
}

func (r *RoutingPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a routing policy that maps each alert severity to an outcome.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the routing policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns this policy.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the routing policy. Unique within your organization.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Human-readable description.",
			},
			"is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "True for the org-wide fallback policy, which cannot be deleted.",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 creation timestamp.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "ISO 8601 last-update timestamp.",
			},
		},
		Blocks: map[string]schema.Block{
			"rule": schema.ListNestedBlock{
				Description: "Exactly four severity rules — one per severity (critical, high, medium, low).",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"severity": schema.StringAttribute{
							Required:    true,
							Description: "Alert severity this rule applies to.",
							Validators: []validator.String{
								stringvalidator.OneOf(validSeverities...),
							},
						},
						"outcome": schema.StringAttribute{
							Required:    true,
							Description: "What happens to alerts of this severity.",
							Validators: []validator.String{
								stringvalidator.OneOf("incident", "provisional_page", "notification", "drop"),
							},
						},
						"escalation_policy_id": schema.StringAttribute{
							Optional:    true,
							Description: "Escalation policy to page for paging outcomes. Null falls back to the alert's component default.",
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeBetween(len(validSeverities), len(validSeverities)),
					&oneRulePerSeverityValidator{},
				},
			},
		},
	}
}

func (r *RoutingPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.ScalingClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected provider data type",
			fmt.Sprintf("Expected *client.ScalingClient, got: %T", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *RoutingPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RoutingPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules := buildRuleInputs(plan.Rules)

	policy, err := r.client.CreateRoutingPolicy(ctx, client.CreateRoutingPolicyRequest{
		Name:        plan.Name.ValueString(),
		Description: stringPtrFromTerraform(plan.Description),
		Rules:       rules,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create routing policy", err.Error())
		return
	}

	r.mapToState(policy, &resp.Diagnostics, &resp.State, ctx)
}

func (r *RoutingPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RoutingPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.client.GetRoutingPolicy(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read routing policy", err.Error())
		return
	}

	r.mapToState(policy, &resp.Diagnostics, &resp.State, ctx)
}

func (r *RoutingPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RoutingPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorState RoutingPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules := buildRuleInputs(plan.Rules)

	policy, err := r.client.UpdateRoutingPolicy(ctx, priorState.ID.ValueString(), client.UpdateRoutingPolicyRequest{
		Name:        plan.Name.ValueString(),
		Description: stringPtrFromTerraform(plan.Description),
		Rules:       rules,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to update routing policy", err.Error())
		return
	}

	r.mapToState(policy, &resp.Diagnostics, &resp.State, ctx)
}

func (r *RoutingPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RoutingPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRoutingPolicy(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Failed to delete routing policy", err.Error())
	}
}

func (r *RoutingPolicyResource) mapToState(policy *client.RoutingPolicyWithRules, diags *diag.Diagnostics, state *tfsdk.State, ctx context.Context) {
	model := RoutingPolicyModel{
		ID:          types.StringValue(policy.ID),
		OrgID:       types.StringValue(policy.OrgID),
		Name:        types.StringValue(policy.Name),
		Description: nullableStringToTerraform(policy.Description),
		IsDefault:   types.BoolValue(policy.IsDefault),
		CreatedAt:   types.StringValue(policy.CreatedAt),
		UpdatedAt:   types.StringValue(policy.UpdatedAt),
		Rules:       make([]RoutingRuleModel, 0, len(policy.Rules)),
	}

	for _, rule := range policy.Rules {
		model.Rules = append(model.Rules, RoutingRuleModel{
			Severity:           types.StringValue(rule.Severity),
			Outcome:            types.StringValue(rule.Outcome),
			EscalationPolicyID: nullableStringToOptionalTerraform(rule.EscalationPolicyID),
		})
	}

	diags.Append(state.Set(ctx, &model)...)
}

func buildRuleInputs(rules []RoutingRuleModel) []client.RoutingRuleInput {
	inputs := make([]client.RoutingRuleInput, len(rules))
	for i, rule := range rules {
		inputs[i] = client.RoutingRuleInput{
			Severity:           rule.Severity.ValueString(),
			Outcome:            rule.Outcome.ValueString(),
			EscalationPolicyID: stringPtrFromTerraform(rule.EscalationPolicyID),
		}
	}
	return inputs
}

func stringPtrFromTerraform(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

func nullableStringToTerraform(s *string) types.String {
	if s == nil {
		return types.StringValue("")
	}
	return types.StringValue(*s)
}

func nullableStringToOptionalTerraform(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
