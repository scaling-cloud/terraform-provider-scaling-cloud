package escalation_policy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	_ resource.Resource              = &EscalationPolicyResource{}
	_ resource.ResourceWithConfigure = &EscalationPolicyResource{}
)

type EscalationPolicyResource struct {
	client *client.ScalingClient
}

type EscalationPolicyModel struct {
	ID          types.String          `tfsdk:"id"`
	OrgID       types.String          `tfsdk:"org_id"`
	Name        types.String          `tfsdk:"name"`
	Description types.String          `tfsdk:"description"`
	CreatedAt   types.String          `tfsdk:"created_at"`
	UpdatedAt   types.String          `tfsdk:"updated_at"`
	Steps       []EscalationStepModel `tfsdk:"step"`
}

type EscalationStepModel struct {
	ID                   types.String `tfsdk:"id"`
	Position             types.Int64  `tfsdk:"position"`
	TargetType           types.String `tfsdk:"target_type"`
	TargetID             types.String `tfsdk:"target_id"`
	EscalateAfterSeconds types.Int64  `tfsdk:"escalate_after_seconds"`
	CreatedAt            types.String `tfsdk:"created_at"`
	UpdatedAt            types.String `tfsdk:"updated_at"`
}

func NewResource() resource.Resource {
	return &EscalationPolicyResource{}
}

func (r *EscalationPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_escalation_policy"
}

func (r *EscalationPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an escalation policy with ordered steps.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the escalation policy.",
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
				Description: "Display name of the escalation policy.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Human-readable description.",
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
			"step": schema.ListNestedBlock{
				Description: "Ordered escalation steps. Position is determined by list order.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the step.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"position": schema.Int64Attribute{
							Computed:    true,
							Description: "1-based position, computed from block order.",
						},
						"target_type": schema.StringAttribute{
							Required:    true,
							Description: "Type of the escalation target.",
							Validators: []validator.String{
								stringvalidator.OneOf("schedule"),
							},
						},
						"target_id": schema.StringAttribute{
							Required:    true,
							Description: "ID of the target on-call schedule.",
						},
						"escalate_after_seconds": schema.Int64Attribute{
							Required:    true,
							Description: "Seconds to wait before escalating to next step. Minimum 60.",
							Validators: []validator.Int64{
								int64validator.AtLeast(60),
							},
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *EscalationPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EscalationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EscalationPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	steps := buildStepInputs(plan.Steps)

	policy, err := r.client.CreateEscalationPolicy(ctx, client.CreateEscalationPolicyRequest{
		Name:        plan.Name.ValueString(),
		Description: stringPtrFromTerraform(plan.Description),
		Steps:       steps,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create escalation policy", err.Error())
		return
	}

	r.mapToState(policy, &resp.Diagnostics, &resp.State, ctx)
}

func (r *EscalationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EscalationPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policy, err := r.client.GetEscalationPolicy(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read escalation policy", err.Error())
		return
	}

	r.mapToState(policy, &resp.Diagnostics, &resp.State, ctx)
}

func (r *EscalationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EscalationPolicyModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorState EscalationPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	steps := buildStepInputs(plan.Steps)

	policy, err := r.client.UpdateEscalationPolicy(ctx, priorState.ID.ValueString(), client.UpdateEscalationPolicyRequest{
		Name:        plan.Name.ValueString(),
		Description: stringPtrFromTerraform(plan.Description),
		Steps:       steps,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to update escalation policy", err.Error())
		return
	}

	r.mapToState(policy, &resp.Diagnostics, &resp.State, ctx)
}

func (r *EscalationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state EscalationPolicyModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEscalationPolicy(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Failed to delete escalation policy", err.Error())
	}
}

func (r *EscalationPolicyResource) mapToState(policy *client.EscalationPolicyWithSteps, diags *diag.Diagnostics, state *tfsdk.State, ctx context.Context) {
	model := EscalationPolicyModel{
		ID:          types.StringValue(policy.ID),
		OrgID:       types.StringValue(policy.OrgID),
		Name:        types.StringValue(policy.Name),
		Description: nullableStringToTerraform(policy.Description),
		CreatedAt:   types.StringValue(policy.CreatedAt),
		UpdatedAt:   types.StringValue(policy.UpdatedAt),
		Steps:       make([]EscalationStepModel, 0, len(policy.Steps)),
	}

	for _, s := range policy.Steps {
		model.Steps = append(model.Steps, EscalationStepModel{
			ID:                   types.StringValue(s.ID),
			Position:             types.Int64Value(int64(s.Position)),
			TargetType:           types.StringValue(s.TargetType),
			TargetID:             types.StringValue(s.TargetID),
			EscalateAfterSeconds: types.Int64Value(int64(s.EscalateAfterSeconds)),
			CreatedAt:            types.StringValue(s.CreatedAt),
			UpdatedAt:            types.StringValue(s.UpdatedAt),
		})
	}

	diags.Append(state.Set(ctx, &model)...)
}

func buildStepInputs(steps []EscalationStepModel) []client.EscalationStepInput {
	inputs := make([]client.EscalationStepInput, len(steps))
	for i, s := range steps {
		inputs[i] = client.EscalationStepInput{
			Position:             i + 1,
			TargetType:           s.TargetType.ValueString(),
			TargetID:             s.TargetID.ValueString(),
			EscalateAfterSeconds: int(s.EscalateAfterSeconds.ValueInt64()),
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
