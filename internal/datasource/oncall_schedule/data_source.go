package oncall_schedule

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
	_ datasource.DataSource              = &OncallScheduleDataSource{}
	_ datasource.DataSourceWithConfigure = &OncallScheduleDataSource{}
)

type OncallScheduleDataSource struct {
	client *client.ScalingClient
}

type OncallScheduleModel struct {
	ID          types.String `tfsdk:"id"`
	OrgID       types.String `tfsdk:"org_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Timezone    types.String `tfsdk:"timezone"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func NewDataSource() datasource.DataSource {
	return &OncallScheduleDataSource{}
}

func (d *OncallScheduleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oncall_schedule"
}

func (d *OncallScheduleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an on-call schedule by name and exposes its id, so an escalation step can target " +
			"a schedule created out-of-band without hardcoding the generated id.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the on-call schedule to look up. Must be unique within the organization.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the matched on-call schedule.",
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns the schedule.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable description, or null if unset.",
			},
			"timezone": schema.StringAttribute{
				Computed:    true,
				Description: "IANA timezone the schedule's rotations are anchored to.",
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

func (d *OncallScheduleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OncallScheduleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config OncallScheduleModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schedules, err := d.client.ListOncallSchedules(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list on-call schedules", err.Error())
		return
	}

	name := config.Name.ValueString()
	match, err := lookup.One(schedules, "on-call schedule", "name", name,
		func(s client.OncallSchedule) string { return s.Name })
	if err != nil {
		resp.Diagnostics.AddError("On-call schedule lookup failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, mapToModel(&match))...)
}

func mapToModel(s *client.OncallSchedule) OncallScheduleModel {
	return OncallScheduleModel{
		ID:          types.StringValue(s.ID),
		OrgID:       types.StringValue(s.OrgID),
		Name:        types.StringValue(s.Name),
		Description: nullableString(s.Description),
		Timezone:    types.StringValue(s.Timezone),
		CreatedAt:   types.StringValue(s.CreatedAt),
		UpdatedAt:   types.StringValue(s.UpdatedAt),
	}
}

func nullableString(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
