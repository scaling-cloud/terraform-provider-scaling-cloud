package working_hours

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
	_ datasource.DataSource              = &WorkingHoursDataSource{}
	_ datasource.DataSourceWithConfigure = &WorkingHoursDataSource{}
)

type WorkingHoursDataSource struct {
	client *client.ScalingClient
}

type WorkingHoursModel struct {
	ID        types.String `tfsdk:"id"`
	OrgID     types.String `tfsdk:"org_id"`
	Name      types.String `tfsdk:"name"`
	Timezone  types.String `tfsdk:"timezone"`
	CreatedAt types.String `tfsdk:"created_at"`
	UpdatedAt types.String `tfsdk:"updated_at"`
}

func NewDataSource() datasource.DataSource {
	return &WorkingHoursDataSource{}
}

func (d *WorkingHoursDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_working_hours"
}

func (d *WorkingHoursDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up a working-hours set by name and exposes its id, so an escalation step condition " +
			"can reference a set created out-of-band without hardcoding the generated id.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the working-hours set to look up. Must be unique within the organization.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the matched working-hours set.",
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns the working-hours set.",
			},
			"timezone": schema.StringAttribute{
				Computed:    true,
				Description: "IANA timezone the windows are evaluated in.",
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

func (d *WorkingHoursDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *WorkingHoursDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config WorkingHoursModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sets, err := d.client.ListWorkingHours(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list working-hours sets", err.Error())
		return
	}

	name := config.Name.ValueString()
	match, err := lookup.One(sets, "working-hours set", "name", name,
		func(w client.WorkingHours) string { return w.Name })
	if err != nil {
		resp.Diagnostics.AddError("Working-hours lookup failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, mapToModel(&match))...)
}

func mapToModel(w *client.WorkingHours) WorkingHoursModel {
	return WorkingHoursModel{
		ID:        types.StringValue(w.ID),
		OrgID:     types.StringValue(w.OrgID),
		Name:      types.StringValue(w.Name),
		Timezone:  types.StringValue(w.Timezone),
		CreatedAt: types.StringValue(w.CreatedAt),
		UpdatedAt: types.StringValue(w.UpdatedAt),
	}
}
