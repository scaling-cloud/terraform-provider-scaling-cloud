package inbound_integration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
)

var (
	_ resource.Resource                = &InboundIntegrationResource{}
	_ resource.ResourceWithConfigure   = &InboundIntegrationResource{}
	_ resource.ResourceWithImportState = &InboundIntegrationResource{}
)

type InboundIntegrationResource struct {
	client *client.ScalingClient
}

type InboundIntegrationModel struct {
	ID              types.String    `tfsdk:"id"`
	IntegrationID   types.String    `tfsdk:"integration_id"`
	OrgID           types.String    `tfsdk:"org_id"`
	Name            types.String    `tfsdk:"name"`
	ComponentID     types.String    `tfsdk:"component_id"`
	RoutingPolicyID types.String    `tfsdk:"routing_policy_id"`
	Selectors       []SelectorModel `tfsdk:"selector"`
	CreatedAt       types.String    `tfsdk:"created_at"`
	UpdatedAt       types.String    `tfsdk:"updated_at"`
}

type SelectorModel struct {
	RoutingPolicyID types.String   `tfsdk:"routing_policy_id"`
	Matchers        []MatcherModel `tfsdk:"matcher"`
}

type MatcherModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func NewResource() resource.Resource {
	return &InboundIntegrationResource{}
}

func (r *InboundIntegrationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inbound_integration"
}

func (r *InboundIntegrationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the ordered Routing Selectors on an existing inbound integration. " +
			"The integration itself is provisioned out-of-band (for example by connecting a source in the UI); " +
			"this resource owns only its selectors, which are replaced wholesale on every apply.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the inbound integration (equal to integration_id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"integration_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the existing inbound integration whose selectors this resource manages.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns this integration.",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Display name of the inbound integration.",
			},
			"component_id": schema.StringAttribute{
				Computed:    true,
				Description: "Component this integration's alerts resolve to.",
			},
			"routing_policy_id": schema.StringAttribute{
				Computed:    true,
				Description: "Routing policy override applied when no selector row matches. Null means the org default policy.",
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
			"selector": schema.ListNestedBlock{
				Description: "Ordered Routing Selector rows. The first row whose matchers all match an alert chooses its routing_policy_id; a miss falls through to routing_policy_id or the org default. Order is significant.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"routing_policy_id": schema.StringAttribute{
							Required:    true,
							Description: "Routing policy this row selects when its matchers all match.",
						},
					},
					Blocks: map[string]schema.Block{
						"matcher": schema.ListNestedBlock{
							Description: "Attribute equality matchers, ANDed together. Severity is never a matcher.",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Required:    true,
										Description: "Declared attribute key to match on.",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 255),
										},
									},
									"value": schema.StringAttribute{
										Required:    true,
										Description: "Value the attribute must equal (case-insensitively).",
										Validators: []validator.String{
											stringvalidator.LengthBetween(1, 255),
										},
									},
								},
							},
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
						},
					},
				},
			},
		},
	}
}

func (r *InboundIntegrationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *InboundIntegrationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("integration_id"), req, resp)
}

func (r *InboundIntegrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan InboundIntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.client.SetInboundSelectors(ctx, plan.IntegrationID.ValueString(), client.SetSelectorsRequest{
		Selectors: selectorsToInput(plan.Selectors),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to set inbound integration selectors", err.Error())
		return
	}

	r.mapToState(ctx, integration, &resp.Diagnostics, &resp.State)
}

func (r *InboundIntegrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state InboundIntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.client.GetInboundIntegration(ctx, state.IntegrationID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read inbound integration", err.Error())
		return
	}

	r.mapToState(ctx, integration, &resp.Diagnostics, &resp.State)
}

func (r *InboundIntegrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan InboundIntegrationModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integration, err := r.client.SetInboundSelectors(ctx, plan.IntegrationID.ValueString(), client.SetSelectorsRequest{
		Selectors: selectorsToInput(plan.Selectors),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to set inbound integration selectors", err.Error())
		return
	}

	r.mapToState(ctx, integration, &resp.Diagnostics, &resp.State)
}

// Delete relinquishes management by clearing the selectors this resource owns.
// The inbound integration itself has no public delete; it lives on past the
// resource's lifecycle.
func (r *InboundIntegrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state InboundIntegrationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.SetInboundSelectors(ctx, state.IntegrationID.ValueString(), client.SetSelectorsRequest{
		Selectors: []client.RoutingSelector{},
	})
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Failed to clear inbound integration selectors", err.Error())
	}
}

func (r *InboundIntegrationResource) mapToState(ctx context.Context, integration *client.InboundIntegration, diags *diag.Diagnostics, state *tfsdk.State) {
	model := InboundIntegrationModel{
		ID:              types.StringValue(integration.ID),
		IntegrationID:   types.StringValue(integration.ID),
		OrgID:           types.StringValue(integration.OrgID),
		Name:            types.StringValue(integration.Name),
		ComponentID:     types.StringValue(integration.ComponentID),
		RoutingPolicyID: nullableStringToOptional(integration.RoutingPolicyID),
		Selectors:       selectorsToState(integration.Selectors),
		CreatedAt:       types.StringValue(integration.CreatedAt),
		UpdatedAt:       types.StringValue(integration.UpdatedAt),
	}

	diags.Append(state.Set(ctx, &model)...)
}

// selectorsToInput converts the configured selector rows into request inputs,
// preserving order exactly. A nil/empty plan yields a non-nil empty slice so the
// PUT clears any prior selectors (full-replacement, ADR-0039).
func selectorsToInput(models []SelectorModel) []client.RoutingSelector {
	out := make([]client.RoutingSelector, 0, len(models))
	for _, m := range models {
		matchers := make([]client.SelectorMatcher, 0, len(m.Matchers))
		for _, mm := range m.Matchers {
			matchers = append(matchers, client.SelectorMatcher{
				Key:   mm.Key.ValueString(),
				Value: mm.Value.ValueString(),
			})
		}
		out = append(out, client.RoutingSelector{
			Matchers:        matchers,
			RoutingPolicyID: m.RoutingPolicyID.ValueString(),
		})
	}
	return out
}

func selectorsToState(selectors []client.RoutingSelector) []SelectorModel {
	if len(selectors) == 0 {
		return nil
	}
	out := make([]SelectorModel, 0, len(selectors))
	for _, s := range selectors {
		matchers := make([]MatcherModel, 0, len(s.Matchers))
		for _, m := range s.Matchers {
			matchers = append(matchers, MatcherModel{
				Key:   types.StringValue(m.Key),
				Value: types.StringValue(m.Value),
			})
		}
		out = append(out, SelectorModel{
			RoutingPolicyID: types.StringValue(s.RoutingPolicyID),
			Matchers:        matchers,
		})
	}
	return out
}

func nullableStringToOptional(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
