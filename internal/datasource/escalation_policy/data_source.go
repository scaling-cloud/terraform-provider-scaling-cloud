package escalation_policy

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
	_ datasource.DataSource              = &EscalationPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &EscalationPolicyDataSource{}
)

type EscalationPolicyDataSource struct {
	client *client.ScalingClient
}

type EscalationPolicyModel struct {
	ID          types.String `tfsdk:"id"`
	OrgID       types.String `tfsdk:"org_id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func NewDataSource() datasource.DataSource {
	return &EscalationPolicyDataSource{}
}

func (d *EscalationPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_escalation_policy"
}

func (d *EscalationPolicyDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up an escalation policy by name and exposes its id, so a routing policy rule " +
			"can reference a policy created out-of-band without hardcoding the generated id.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the escalation policy to look up. Must be unique within the organization.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the matched escalation policy.",
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns the policy.",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Human-readable description, or null if unset.",
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

func (d *EscalationPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EscalationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config EscalationPolicyModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	policies, err := d.client.ListEscalationPolicies(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list escalation policies", err.Error())
		return
	}

	name := config.Name.ValueString()
	match, err := lookup.One(policies, "escalation policy", "name", name,
		func(p client.EscalationPolicy) string { return p.Name })
	if err != nil {
		resp.Diagnostics.AddError("Escalation policy lookup failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, mapToModel(&match))...)
}

func mapToModel(p *client.EscalationPolicy) EscalationPolicyModel {
	return EscalationPolicyModel{
		ID:          types.StringValue(p.ID),
		OrgID:       types.StringValue(p.OrgID),
		Name:        types.StringValue(p.Name),
		Description: nullableString(p.Description),
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
