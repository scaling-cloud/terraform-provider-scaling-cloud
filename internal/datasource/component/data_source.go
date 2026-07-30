package component

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/datasource/lookup"
)

var (
	_ datasource.DataSource              = &ComponentDataSource{}
	_ datasource.DataSourceWithConfigure = &ComponentDataSource{}
)

type ComponentDataSource struct {
	client *client.ScalingClient
}

type ComponentModel struct {
	ID                types.String `tfsdk:"id"`
	OrgID             types.String `tfsdk:"org_id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	Aliases           types.Set    `tfsdk:"aliases"`
	OperationalStatus types.String `tfsdk:"operational_status"`
	CreatedAt         types.String `tfsdk:"created_at"`
	UpdatedAt         types.String `tfsdk:"updated_at"`
}

func NewDataSource() datasource.DataSource {
	return &ComponentDataSource{}
}

func (d *ComponentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_component"
}

func (d *ComponentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up a component by name and exposes its id, so a configuration can reference a " +
			"component created out-of-band without hardcoding the generated id.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the component to look up. Must be unique within the organization.",
			},
			"alias": schema.StringAttribute{
				Optional:    true,
				Description: "Optional alias to narrow the lookup. When set, only components whose aliases " +
					"include this value will be considered. Use this to disambiguate when multiple " +
					"components share a name in the organization.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the matched component.",
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns the component.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable description, or null if unset.",
			},
			"aliases": schema.SetAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Alternate names inbound alerts match against to resolve this component.",
			},
			"operational_status": schema.StringAttribute{
				Computed:    true,
				Description: "Current operational status (for example operational, degraded, major_outage, maintenance).",
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

func (d *ComponentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ComponentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ComponentModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	components, err := d.client.ListComponents(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list components", err.Error())
		return
	}

	var alias string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("alias"), &alias)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if alias != "" {
		components = filterByAlias(components, alias)
		if len(components) == 0 {
			resp.Diagnostics.AddError("Component lookup failed",
				fmt.Sprintf("no component found with alias %q", alias))
			return
		}
	}

	name := config.Name.ValueString()
	match, err := lookup.One(components, "component", "name", name,
		func(c client.Component) string { return c.Name })
	if err != nil {
		resp.Diagnostics.AddError("Component lookup failed", err.Error())
		return
	}

	model, diags := mapToModel(ctx, &match)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func filterByAlias(components []client.Component, alias string) []client.Component {
	var filtered []client.Component
	for _, c := range components {
		for _, a := range c.Aliases {
			if a == alias {
				filtered = append(filtered, c)
				break
			}
		}
	}
	return filtered
}

func mapToModel(ctx context.Context, c *client.Component) (ComponentModel, diag.Diagnostics) {
	aliases := c.Aliases
	if aliases == nil {
		aliases = []string{}
	}
	aliasSet, diags := types.SetValueFrom(ctx, types.StringType, aliases)

	return ComponentModel{
		ID:                types.StringValue(c.ID),
		OrgID:             types.StringValue(c.OrgID),
		Name:              types.StringValue(c.Name),
		Description:       nullableString(c.Description),
		Aliases:           aliasSet,
		OperationalStatus: types.StringValue(c.OperationalStatus),
		CreatedAt:         types.StringValue(c.CreatedAt),
		UpdatedAt:         types.StringValue(c.UpdatedAt),
	}, diags
}

func nullableString(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
