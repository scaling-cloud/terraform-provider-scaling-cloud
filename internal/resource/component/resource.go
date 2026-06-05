package component

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
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
	_ resource.Resource                = &ComponentResource{}
	_ resource.ResourceWithConfigure   = &ComponentResource{}
	_ resource.ResourceWithImportState = &ComponentResource{}
)

type ComponentResource struct {
	client *client.ScalingClient
}

type ComponentModel struct {
	ID          types.String `tfsdk:"id"`
	OrgID       types.String `tfsdk:"org_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Aliases     types.Set    `tfsdk:"aliases"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func NewResource() resource.Resource {
	return &ComponentResource{}
}

func (r *ComponentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component"
}

func (r *ComponentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a component in the service catalog. Components can carry aliases — alternate names that inbound alerts match against to resolve the component.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the component.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns this component.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the component. Unique within your organization.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Human-readable description.",
			},
			"aliases": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Alternate names alerts can resolve this component by. Replaced wholesale on every apply.",
				Validators: []validator.Set{
					setvalidator.SizeAtMost(50),
					setvalidator.ValueStringsAre(
						stringvalidator.LengthBetween(1, 255),
					),
				},
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
	}
}

func (r *ComponentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ComponentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, tfsdkPathID, req, resp)
}

func (r *ComponentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ComponentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	aliases, diags := aliasesToInput(ctx, plan.Aliases)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	component, err := r.client.CreateComponent(ctx, client.CreateComponentRequest{
		Name:        plan.Name.ValueString(),
		Description: stringPtrFromTerraform(plan.Description),
		Aliases:     aliases,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create component", err.Error())
		return
	}

	r.mapToState(ctx, component, &resp.Diagnostics, &resp.State)
}

func (r *ComponentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ComponentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	component, err := r.client.GetComponent(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read component", err.Error())
		return
	}

	r.mapToState(ctx, component, &resp.Diagnostics, &resp.State)
}

func (r *ComponentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ComponentModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorState ComponentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	aliases, diags := aliasesToInput(ctx, plan.Aliases)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	component, err := r.client.UpdateComponent(ctx, priorState.ID.ValueString(), client.UpdateComponentRequest{
		Name:        plan.Name.ValueString(),
		Description: stringPtrFromTerraform(plan.Description),
		Aliases:     aliases,
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to update component", err.Error())
		return
	}

	r.mapToState(ctx, component, &resp.Diagnostics, &resp.State)
}

func (r *ComponentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ComponentModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteComponent(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Failed to delete component", err.Error())
	}
}

func (r *ComponentResource) mapToState(ctx context.Context, component *client.Component, diags *diag.Diagnostics, state *tfsdk.State) {
	aliases, aliasDiags := aliasesToState(ctx, component.Aliases)
	diags.Append(aliasDiags...)
	if diags.HasError() {
		return
	}

	model := ComponentModel{
		ID:          types.StringValue(component.ID),
		OrgID:       types.StringValue(component.OrgID),
		Name:        types.StringValue(component.Name),
		Description: nullableStringToTerraform(component.Description),
		Aliases:     aliases,
		CreatedAt:   types.StringValue(component.CreatedAt),
		UpdatedAt:   types.StringValue(component.UpdatedAt),
	}

	diags.Append(state.Set(ctx, &model)...)
}

// aliasesToInput converts the configured alias set into a sorted, non-nil slice
// for the request body. Sorting keeps the PUT-replacement payload deterministic;
// a null/unknown set becomes an empty slice so the server clears prior aliases.
func aliasesToInput(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	aliases := []string{}
	if set.IsNull() || set.IsUnknown() {
		return aliases, diags
	}
	diags.Append(set.ElementsAs(ctx, &aliases, false)...)
	if diags.HasError() {
		return aliases, diags
	}
	sort.Strings(aliases)
	return aliases, diags
}

// aliasesToState converts the API's alias slice into a known (never null) set so
// an empty alias list round-trips as an empty set rather than null.
func aliasesToState(ctx context.Context, aliases []string) (types.Set, diag.Diagnostics) {
	if aliases == nil {
		aliases = []string{}
	}
	return types.SetValueFrom(ctx, types.StringType, aliases)
}
