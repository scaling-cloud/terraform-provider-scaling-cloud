package inbound_integration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/datasource/lookup"
)

var (
	_ datasource.DataSource              = &InboundIntegrationDataSource{}
	_ datasource.DataSourceWithConfigure = &InboundIntegrationDataSource{}
)

type InboundIntegrationDataSource struct {
	client *client.ScalingClient
}

type InboundIntegrationModel struct {
	ID              types.String `tfsdk:"id"`
	OrgID           types.String `tfsdk:"org_id"`
	Name            types.String `tfsdk:"name"`
	ComponentID     types.String `tfsdk:"component_id"`
	RoutingPolicyID types.String `tfsdk:"routing_policy_id"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

func NewDataSource() datasource.DataSource {
	return &InboundIntegrationDataSource{}
}

func (d *InboundIntegrationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inbound_integration"
}

func (d *InboundIntegrationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an inbound integration by name and exposes its id. Inbound integrations are " +
			"provisioned out-of-band (for example by connecting a source via an OAuth install), so this data " +
			"source lets a configuration reference one — for instance to manage its selectors with the " +
			"scaling_inbound_integration resource — without hardcoding the generated id.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the inbound integration to look up. Must be unique within the organization.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the matched inbound integration.",
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns the integration.",
			},
			"component_id": schema.StringAttribute{
				Computed:    true,
				Description: "Component this integration's alerts resolve to.",
			},
			"routing_policy_id": schema.StringAttribute{
				Computed:    true,
				Description: "Routing policy applied when no selector row matches. Null means the org default policy.",
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

func (d *InboundIntegrationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.client = c
}

func (d *InboundIntegrationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config InboundIntegrationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	integrations, err := d.client.ListInboundIntegrations(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list inbound integrations", err.Error())
		return
	}

	name := config.Name.ValueString()
	match, err := lookup.One(integrations, "inbound integration", "name", name,
		func(i client.InboundIntegration) string { return i.Name })
	if err != nil {
		resp.Diagnostics.AddError("Inbound integration lookup failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, mapToModel(&match))...)
}

func mapToModel(i *client.InboundIntegration) InboundIntegrationModel {
	return InboundIntegrationModel{
		ID:              types.StringValue(i.ID),
		OrgID:           types.StringValue(i.OrgID),
		Name:            types.StringValue(i.Name),
		ComponentID:     types.StringValue(i.ComponentID),
		RoutingPolicyID: nullableString(i.RoutingPolicyID),
		CreatedAt:       types.StringValue(i.CreatedAt),
		UpdatedAt:       types.StringValue(i.UpdatedAt),
	}
}

func nullableString(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
