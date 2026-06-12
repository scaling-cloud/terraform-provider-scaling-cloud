package routing_policy

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
	_ datasource.DataSource              = &RoutingPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &RoutingPolicyDataSource{}
)

type RoutingPolicyDataSource struct {
	client *client.ScalingClient
}

type RoutingPolicyModel struct {
	ID          types.String `tfsdk:"id"`
	OrgID       types.String `tfsdk:"org_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IsDefault   types.Bool   `tfsdk:"is_default"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func NewDataSource() datasource.DataSource {
	return &RoutingPolicyDataSource{}
}

func (d *RoutingPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_routing_policy"
}

func (d *RoutingPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up a routing policy by name and exposes its id, so an inbound integration selector " +
			"can reference a policy created out-of-band without hardcoding the generated id.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the routing policy to look up. Must be unique within the organization.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the matched routing policy.",
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns the policy.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable description, or null if unset.",
			},
			"is_default": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this is the organization's default routing policy.",
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

func (d *RoutingPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *RoutingPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config RoutingPolicyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policies, err := d.client.ListRoutingPolicies(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list routing policies", err.Error())
		return
	}

	name := config.Name.ValueString()
	match, err := lookup.One(policies, "routing policy", "name", name,
		func(p client.RoutingPolicy) string { return p.Name })
	if err != nil {
		resp.Diagnostics.AddError("Routing policy lookup failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, mapToModel(&match))...)
}

func mapToModel(p *client.RoutingPolicy) RoutingPolicyModel {
	return RoutingPolicyModel{
		ID:          types.StringValue(p.ID),
		OrgID:       types.StringValue(p.OrgID),
		Name:        types.StringValue(p.Name),
		Description: nullableString(p.Description),
		IsDefault:   types.BoolValue(p.IsDefault),
		CreatedAt:   types.StringValue(p.CreatedAt),
		UpdatedAt:   types.StringValue(p.UpdatedAt),
	}
}

func nullableString(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
