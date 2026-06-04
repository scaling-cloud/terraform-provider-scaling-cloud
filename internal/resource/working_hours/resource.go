package working_hours

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	_ resource.Resource              = &WorkingHoursResource{}
	_ resource.ResourceWithConfigure = &WorkingHoursResource{}
)

// timeOfDayRegex matches a 24-hour HH:MM string, mirroring the /v1 schema.
var timeOfDayRegex = regexp.MustCompile(`^([01]\d|2[0-3]):[0-5]\d$`)

type WorkingHoursResource struct {
	client *client.ScalingClient
}

type WorkingHoursModel struct {
	ID        types.String  `tfsdk:"id"`
	OrgID     types.String  `tfsdk:"org_id"`
	Name      types.String  `tfsdk:"name"`
	Timezone  types.String  `tfsdk:"timezone"`
	CreatedAt types.String  `tfsdk:"created_at"`
	UpdatedAt types.String  `tfsdk:"updated_at"`
	Windows   []WindowModel `tfsdk:"window"`
}

type WindowModel struct {
	Days  []types.Int64 `tfsdk:"days"`
	Start types.String  `tfsdk:"start"`
	End   types.String  `tfsdk:"end"`
}

func NewResource() resource.Resource {
	return &WorkingHoursResource{}
}

func (r *WorkingHoursResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_working_hours"
}

func (r *WorkingHoursResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a named, reusable Working Hours set — a timezone plus weekly windows — used by escalation policy step conditions for follow-the-sun routing.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the working hours set.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns this working hours set.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the working hours set. Unique within your organization.",
			},
			"timezone": schema.StringAttribute{
				Required:    true,
				Description: "IANA timezone the windows are interpreted in (e.g. \"Europe/London\"). The timezone lives on the set; it is never inferred from a targeted schedule.",
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
			"window": schema.ListNestedBlock{
				Description: "Weekly windows during which these hours are active. At least one window is required.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"days": schema.ListAttribute{
							Required:    true,
							ElementType: types.Int64Type,
							Description: "ISO weekday numbers this window applies to (1=Monday through 7=Sunday).",
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
						},
						"start": schema.StringAttribute{
							Required:    true,
							Description: "Window start time as 24-hour HH:MM in the set's timezone.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(timeOfDayRegex, "must be a 24-hour HH:MM time"),
							},
						},
						"end": schema.StringAttribute{
							Required:    true,
							Description: "Window end time as 24-hour HH:MM in the set's timezone.",
							Validators: []validator.String{
								stringvalidator.RegexMatches(timeOfDayRegex, "must be a 24-hour HH:MM time"),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (r *WorkingHoursResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WorkingHoursResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WorkingHoursModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.CreateWorkingHours(ctx, client.CreateWorkingHoursRequest{
		Name:     plan.Name.ValueString(),
		Timezone: plan.Timezone.ValueString(),
		Windows:  buildWindowInputs(plan.Windows),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create working hours", err.Error())
		return
	}

	r.mapToState(set, &resp.Diagnostics, &resp.State, ctx)
}

func (r *WorkingHoursResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WorkingHoursModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.GetWorkingHours(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read working hours", err.Error())
		return
	}

	r.mapToState(set, &resp.Diagnostics, &resp.State, ctx)
}

func (r *WorkingHoursResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WorkingHoursModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorState WorkingHoursModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	set, err := r.client.UpdateWorkingHours(ctx, priorState.ID.ValueString(), client.UpdateWorkingHoursRequest{
		Name:     plan.Name.ValueString(),
		Timezone: plan.Timezone.ValueString(),
		Windows:  buildWindowInputs(plan.Windows),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to update working hours", err.Error())
		return
	}

	r.mapToState(set, &resp.Diagnostics, &resp.State, ctx)
}

func (r *WorkingHoursResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WorkingHoursModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWorkingHours(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Failed to delete working hours", err.Error())
	}
}

func (r *WorkingHoursResource) mapToState(set *client.WorkingHours, diags *diag.Diagnostics, state *tfsdk.State, ctx context.Context) {
	model := WorkingHoursModel{
		ID:        types.StringValue(set.ID),
		OrgID:     types.StringValue(set.OrgID),
		Name:      types.StringValue(set.Name),
		Timezone:  types.StringValue(set.Timezone),
		CreatedAt: types.StringValue(set.CreatedAt),
		UpdatedAt: types.StringValue(set.UpdatedAt),
		Windows:   make([]WindowModel, 0, len(set.Windows)),
	}

	for _, w := range set.Windows {
		days := make([]types.Int64, 0, len(w.Days))
		for _, d := range w.Days {
			days = append(days, types.Int64Value(int64(d)))
		}
		model.Windows = append(model.Windows, WindowModel{
			Days:  days,
			Start: types.StringValue(w.Start),
			End:   types.StringValue(w.End),
		})
	}

	diags.Append(state.Set(ctx, &model)...)
}

func buildWindowInputs(windows []WindowModel) []client.WorkingHoursWindow {
	inputs := make([]client.WorkingHoursWindow, len(windows))
	for i, w := range windows {
		days := make([]int, 0, len(w.Days))
		for _, d := range w.Days {
			days = append(days, int(d.ValueInt64()))
		}
		inputs[i] = client.WorkingHoursWindow{
			Days:  days,
			Start: w.Start.ValueString(),
			End:   w.End.ValueString(),
		}
	}
	return inputs
}
