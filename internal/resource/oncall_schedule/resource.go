package oncall_schedule

import (
	"context"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
	_ resource.Resource              = &OncallScheduleResource{}
	_ resource.ResourceWithConfigure = &OncallScheduleResource{}
)

type OncallScheduleResource struct {
	client *client.ScalingClient
}

type OncallScheduleModel struct {
	ID          types.String       `tfsdk:"id"`
	OrgID       types.String       `tfsdk:"org_id"`
	Name        types.String       `tfsdk:"name"`
	Description types.String       `tfsdk:"description"`
	Timezone    types.String       `tfsdk:"timezone"`
	CreatedAt   types.String       `tfsdk:"created_at"`
	UpdatedAt   types.String       `tfsdk:"updated_at"`
	Layers      []OncallLayerModel `tfsdk:"layer"`
}

type OncallLayerModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	RotationType       types.String `tfsdk:"rotation_type"`
	RotationLengthDays types.Int64  `tfsdk:"rotation_length_days"`
	HandoffTime        types.String `tfsdk:"handoff_time"`
	EffectiveFrom      types.String `tfsdk:"effective_from"`
	EffectiveUntil     types.String `tfsdk:"effective_until"`
	ParticipantIDs     types.List   `tfsdk:"participant_ids"`
	CreatedAt          types.String `tfsdk:"created_at"`
	UpdatedAt          types.String `tfsdk:"updated_at"`
}

func NewResource() resource.Resource {
	return &OncallScheduleResource{}
}

func (r *OncallScheduleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oncall_schedule"
}

func (r *OncallScheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an on-call schedule with rotation layers.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for the schedule.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_id": schema.StringAttribute{
				Computed:    true,
				Description: "Organization that owns this schedule.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Display name of the schedule.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
				Description: "Human-readable description.",
			},
			"timezone": schema.StringAttribute{
				Required:    true,
				Description: "IANA timezone for the schedule (e.g. \"America/New_York\").",
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
			"layer": schema.ListNestedBlock{
				Description: "Rotation layers within this schedule.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Unique identifier for the layer.",
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"name": schema.StringAttribute{
							Required:    true,
							Description: "Display name of the layer.",
						},
						"rotation_type": schema.StringAttribute{
							Required:    true,
							Description: "Rotation type: \"daily\", \"weekly\", or \"custom\".",
							Validators: []validator.String{
								stringvalidator.OneOf("daily", "weekly", "custom"),
							},
						},
						"rotation_length_days": schema.Int64Attribute{
							Required:    true,
							Description: "Length of each rotation in days.",
							Validators: []validator.Int64{
								int64validator.AtLeast(1),
							},
						},
						"handoff_time": schema.StringAttribute{
							Required:    true,
							Description: "Time of day for rotation handoff in HH:MM format.",
						},
						"effective_from": schema.StringAttribute{
							Required:    true,
							Description: "ISO 8601 datetime when this layer becomes active.",
						},
						"effective_until": schema.StringAttribute{
							Optional:    true,
							Description: "ISO 8601 datetime when this layer expires.",
						},
						"participant_ids": schema.ListAttribute{
							Required:    true,
							ElementType: types.StringType,
							Description: "Ordered list of participant user IDs for the rotation.",
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
						},
						"created_at": schema.StringAttribute{
							Computed: true,
						},
						"updated_at": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (r *OncallScheduleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OncallScheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OncallScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	schedule, err := r.client.CreateOncallSchedule(ctx, client.CreateScheduleRequest{
		Name:        plan.Name.ValueString(),
		Description: stringPtrFromTerraform(plan.Description),
		Timezone:    plan.Timezone.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create on-call schedule", err.Error())
		return
	}

	for _, layer := range plan.Layers {
		var participantIDs []string
		resp.Diagnostics.Append(layer.ParticipantIDs.ElementsAs(ctx, &participantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
		_, err := r.client.CreateOncallLayer(ctx, schedule.ID, client.CreateLayerRequest{
			Name:               layer.Name.ValueString(),
			RotationType:       layer.RotationType.ValueString(),
			RotationLengthDays: int(layer.RotationLengthDays.ValueInt64()),
			HandoffTime:        layer.HandoffTime.ValueString(),
			EffectiveFrom:      layer.EffectiveFrom.ValueString(),
			EffectiveUntil:     stringPtrFromTerraform(layer.EffectiveUntil),
			ParticipantIDs:     participantIDs,
		})
		if err != nil {
			resp.Diagnostics.AddError("Failed to create on-call layer", err.Error())
			break
		}
	}

	r.readIntoState(ctx, schedule.ID, &resp.Diagnostics, &resp.State)
}

func (r *OncallScheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OncallScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.readIntoState(ctx, state.ID.ValueString(), &resp.Diagnostics, &resp.State)
}

func (r *OncallScheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OncallScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var priorState OncallScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &priorState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scheduleID := priorState.ID.ValueString()

	_, err := r.client.UpdateOncallSchedule(ctx, scheduleID, client.UpdateScheduleRequest{
		Name:        plan.Name.ValueString(),
		Description: stringPtrFromTerraform(plan.Description),
		Timezone:    plan.Timezone.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to update on-call schedule", err.Error())
		return
	}

	existingByID := make(map[string]bool)
	for _, l := range priorState.Layers {
		if !l.ID.IsNull() && !l.ID.IsUnknown() {
			existingByID[l.ID.ValueString()] = true
		}
	}

	plannedIDs := make(map[string]bool)
	for _, layer := range plan.Layers {
		var participantIDs []string
		resp.Diagnostics.Append(layer.ParticipantIDs.ElementsAs(ctx, &participantIDs, false)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !layer.ID.IsNull() && !layer.ID.IsUnknown() {
			layerID := layer.ID.ValueString()
			plannedIDs[layerID] = true
			_, err := r.client.UpdateOncallLayer(ctx, scheduleID, layerID, client.UpdateLayerRequest{
				Name:               layer.Name.ValueString(),
				RotationType:       layer.RotationType.ValueString(),
				RotationLengthDays: int(layer.RotationLengthDays.ValueInt64()),
				HandoffTime:        layer.HandoffTime.ValueString(),
				EffectiveFrom:      layer.EffectiveFrom.ValueString(),
				EffectiveUntil:     stringPtrFromTerraform(layer.EffectiveUntil),
				ParticipantIDs:     participantIDs,
			})
			if err != nil {
				resp.Diagnostics.AddError("Failed to update on-call layer", err.Error())
				break
			}
		} else {
			_, err := r.client.CreateOncallLayer(ctx, scheduleID, client.CreateLayerRequest{
				Name:               layer.Name.ValueString(),
				RotationType:       layer.RotationType.ValueString(),
				RotationLengthDays: int(layer.RotationLengthDays.ValueInt64()),
				HandoffTime:        layer.HandoffTime.ValueString(),
				EffectiveFrom:      layer.EffectiveFrom.ValueString(),
				EffectiveUntil:     stringPtrFromTerraform(layer.EffectiveUntil),
				ParticipantIDs:     participantIDs,
			})
			if err != nil {
				resp.Diagnostics.AddError("Failed to create on-call layer", err.Error())
				break
			}
		}
	}

	for id := range existingByID {
		if !plannedIDs[id] {
			if err := r.client.DeleteOncallLayer(ctx, scheduleID, id); err != nil && !client.IsNotFound(err) {
				resp.Diagnostics.AddError("Failed to delete on-call layer", err.Error())
				break
			}
		}
	}

	r.readIntoState(ctx, scheduleID, &resp.Diagnostics, &resp.State)
}

func (r *OncallScheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OncallScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOncallSchedule(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Failed to delete on-call schedule", err.Error())
	}
}

func (r *OncallScheduleResource) readIntoState(ctx context.Context, scheduleID string, diags *diag.Diagnostics, state *tfsdk.State) {
	schedule, err := r.client.GetOncallSchedule(ctx, scheduleID)
	if err != nil {
		if client.IsNotFound(err) {
			state.RemoveResource(ctx)
			return
		}
		diags.AddError("Failed to read on-call schedule", err.Error())
		return
	}

	sort.Slice(schedule.Layers, func(i, j int) bool {
		return schedule.Layers[i].ID < schedule.Layers[j].ID
	})

	model := OncallScheduleModel{
		ID:          types.StringValue(schedule.ID),
		OrgID:       types.StringValue(schedule.OrgID),
		Name:        types.StringValue(schedule.Name),
		Description: nullableStringToTerraform(schedule.Description),
		Timezone:    types.StringValue(schedule.Timezone),
		CreatedAt:   types.StringValue(schedule.CreatedAt),
		UpdatedAt:   types.StringValue(schedule.UpdatedAt),
		Layers:      make([]OncallLayerModel, 0, len(schedule.Layers)),
	}

	for _, l := range schedule.Layers {
		pIDs, listDiags := types.ListValueFrom(ctx, types.StringType, l.ParticipantIDs)
		diags.Append(listDiags...)

		model.Layers = append(model.Layers, OncallLayerModel{
			ID:                 types.StringValue(l.ID),
			Name:               types.StringValue(l.Name),
			RotationType:       types.StringValue(l.RotationType),
			RotationLengthDays: types.Int64Value(int64(l.RotationLengthDays)),
			HandoffTime:        types.StringValue(l.HandoffTime),
			EffectiveFrom:      types.StringValue(l.EffectiveFrom),
			EffectiveUntil:     nullableStringToTerraform(l.EffectiveUntil),
			ParticipantIDs:     pIDs,
			CreatedAt:          types.StringValue(l.CreatedAt),
			UpdatedAt:          types.StringValue(l.UpdatedAt),
		})
	}

	diags.Append(state.Set(ctx, &model)...)
}

func stringPtrFromTerraform(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

func nullableStringToTerraform(s *string) types.String {
	if s == nil {
		return types.StringValue("")
	}
	return types.StringValue(*s)
}
