package user

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
	_ datasource.DataSource              = &UserDataSource{}
	_ datasource.DataSourceWithConfigure = &UserDataSource{}
)

type UserDataSource struct {
	client *client.ScalingClient
}

type UserModel struct {
	ID        types.String `tfsdk:"id"`
	Email     types.String `tfsdk:"email"`
	FirstName types.String `tfsdk:"first_name"`
	LastName  types.String `tfsdk:"last_name"`
}

func NewDataSource() datasource.DataSource {
	return &UserDataSource{}
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Looks up a user by email and exposes their id, so a configuration can resolve a user " +
			"to an id — for example to populate an on-call schedule layer's participant_ids.",
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				Required:    true,
				Description: "Email address of the user to look up. Must match exactly one user in the organization.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier of the matched user.",
			},
			"first_name": schema.StringAttribute{
				Computed:    true,
				Description: "The user's first name, or null if unknown.",
			},
			"last_name": schema.StringAttribute{
				Computed:    true,
				Description: "The user's last name, or null if unknown.",
			},
		},
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	users, err := d.client.ListUsers(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list users", err.Error())
		return
	}

	email := config.Email.ValueString()
	match, err := lookup.One(users, "user", "email", email, emailOf)
	if err != nil {
		resp.Diagnostics.AddError("User lookup failed", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, mapToModel(&match))...)
}

// emailOf returns a user's email for matching. A user with no email can never
// match a requested email, so it yields the empty string.
func emailOf(u client.User) string {
	if u.Email == nil {
		return ""
	}
	return *u.Email
}

func mapToModel(u *client.User) UserModel {
	return UserModel{
		ID:        types.StringValue(u.ID),
		Email:     nullableString(u.Email),
		FirstName: nullableString(u.FirstName),
		LastName:  nullableString(u.LastName),
	}
}

func nullableString(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
