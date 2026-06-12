package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/client"
	escalation_policy_data "github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/datasource/escalation_policy"
	inbound_integration_data "github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/datasource/inbound_integration"
	oncall_schedule_data "github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/datasource/oncall_schedule"
	routing_policy_data "github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/datasource/routing_policy"
	working_hours_data "github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/datasource/working_hours"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/resource/component"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/resource/escalation_policy"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/resource/inbound_integration"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/resource/oncall_schedule"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/resource/routing_policy"
	"github.com/scaling-cloud/terraform-provider-scaling-cloud/internal/resource/working_hours"
)

var _ provider.Provider = &ScalingCloudProvider{}

type ScalingCloudProvider struct {
	version string
}

type ScalingCloudProviderModel struct {
	APIKey  types.String `tfsdk:"api_key"`
	BaseURL types.String `tfsdk:"base_url"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ScalingCloudProvider{
			version: version,
		}
	}
}

func (p *ScalingCloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "scaling"
	resp.Version = p.version
}

func (p *ScalingCloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform provider for the Scaling Cloud incident management platform.",
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "API key for authenticating with the Scaling Cloud API. " +
					"Can also be set via the SCALING_CLOUD_API_KEY environment variable.",
				Optional:  true,
				Sensitive: true,
			},
			"base_url": schema.StringAttribute{
				Description: "Base URL of the Scaling Cloud API. " +
					"Defaults to https://api.scaling.cloud. " +
					"Can also be set via the SCALING_CLOUD_BASE_URL environment variable.",
				Optional: true,
			},
		},
	}
}

func (p *ScalingCloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config ScalingCloudProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := resolveStringAttribute(config.APIKey, "SCALING_CLOUD_API_KEY", "")
	baseURL := resolveStringAttribute(config.BaseURL, "SCALING_CLOUD_BASE_URL", "https://api.scaling.cloud")

	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"The provider requires an api_key attribute or SCALING_CLOUD_API_KEY environment variable.",
		)
		return
	}

	c, err := client.NewScalingClient(baseURL, apiKey)
	if err != nil {
		resp.Diagnostics.AddError("Invalid provider configuration", err.Error())
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *ScalingCloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		oncall_schedule.NewResource,
		escalation_policy.NewResource,
		routing_policy.NewResource,
		working_hours.NewResource,
		component.NewResource,
		inbound_integration.NewResource,
	}
}

func (p *ScalingCloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		inbound_integration_data.NewDataSource,
		escalation_policy_data.NewDataSource,
		routing_policy_data.NewDataSource,
		oncall_schedule_data.NewDataSource,
		working_hours_data.NewDataSource,
	}
}

func resolveStringAttribute(attr types.String, envVar, defaultValue string) string {
	if !attr.IsNull() && !attr.IsUnknown() {
		return attr.ValueString()
	}
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	return defaultValue
}
