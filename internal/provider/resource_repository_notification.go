package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lannuttia/terraform-provider-quay/internal/client"
)

var _ resource.Resource = &RepositoryNotificationResource{}
var _ resource.ResourceWithImportState = &RepositoryNotificationResource{}

type RepositoryNotificationResource struct {
	client *client.Client
}

type RepositoryNotificationResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Repository       types.String `tfsdk:"repository"`
	UUID             types.String `tfsdk:"uuid"`
	Title            types.String `tfsdk:"title"`
	Event            types.String `tfsdk:"event"`
	EventConfig      types.Map    `tfsdk:"event_config"`
	Method           types.String `tfsdk:"method"`
	Config           types.Map    `tfsdk:"config"`
	NumberOfFailures types.Int64  `tfsdk:"number_of_failures"`
}

func NewRepositoryNotificationResource() resource.Resource {
	return &RepositoryNotificationResource{}
}

func (r *RepositoryNotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_notification"
}

func (r *RepositoryNotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a notification (webhook) on a Quay repository. Because the Quay API has no update endpoint for notifications, any change to notification attributes will cause the existing notification to be deleted and a new one created.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier (repository/uuid).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"repository": schema.StringAttribute{
				Required:    true,
				Description: "The repository this notification belongs to, in namespace/name format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The server-assigned UUID of the notification.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A human-readable title for the notification.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"event": schema.StringAttribute{
				Required:    true,
				Description: "The event that triggers this notification. One of: repo_push, build_queued, build_start, build_success, build_failure, build_cancelled, vulnerability_found, repo_mirror_sync_started, repo_mirror_sync_success, repo_mirror_sync_failed.",
				Validators: []validator.String{
					stringvalidator.OneOf(
						"repo_push",
						"build_queued",
						"build_start",
						"build_success",
						"build_failure",
						"build_cancelled",
						"vulnerability_found",
						"repo_mirror_sync_started",
						"repo_mirror_sync_success",
						"repo_mirror_sync_failed",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"event_config": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
				Description: "Event-specific configuration. For vulnerability_found, set 'level' to a severity string (e.g. 'High'). For build events, optionally set 'ref-regex'.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"method": schema.StringAttribute{
				Required:    true,
				Description: "The notification delivery method. One of: webhook, email, slack, hipchat, flowdock, quay_notification.",
				Validators: []validator.String{
					stringvalidator.OneOf("webhook", "email", "slack", "hipchat", "flowdock", "quay_notification"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.MapAttribute{
				ElementType: types.StringType,
				Required:    true,
				Description: "Method-specific configuration. For webhook/slack: set 'url'. For email: set 'email'. For hipchat: set 'notification_token' and 'room_id'. For flowdock: set 'flow_api_token'.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"number_of_failures": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of consecutive delivery failures.",
			},
		},
	}
}

func (r *RepositoryNotificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *RepositoryNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RepositoryNotificationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := mapToStringAny(ctx, plan.Config)
	if err != nil {
		resp.Diagnostics.AddError("Error Processing Config", err.Error())
		return
	}
	eventConfig, err := mapToStringAny(ctx, plan.EventConfig)
	if err != nil {
		resp.Diagnostics.AddError("Error Processing Event Config", err.Error())
		return
	}

	notif, err := r.client.CreateRepoNotification(ctx, plan.Repository.ValueString(), client.CreateNotificationRequest{
		Event:       plan.Event.ValueString(),
		Method:      plan.Method.ValueString(),
		Config:      config,
		EventConfig: eventConfig,
		Title:       plan.Title.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Repository Notification", err.Error())
		return
	}

	d := r.mapToState(ctx, notif, plan.Repository.ValueString(), &plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RepositoryNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RepositoryNotificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	notif, err := r.client.GetRepoNotification(ctx, state.Repository.ValueString(), state.UUID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Repository Notification", err.Error())
		return
	}

	d := r.mapToState(ctx, notif, state.Repository.ValueString(), &state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RepositoryNotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID: namespace/repo_name/uuid  (3 parts, last part is the UUID)
	// We use SplitN with n=3 so the UUID stays intact as the third element.
	parts, diags := splitImportID(resp, req.ID, 3, "my-org/my-repo/32bb0ef6-ad69-4ccb-b64e-ddbc8bfd444d")
	if diags.HasError() {
		return
	}
	repository := parts[0] + "/" + parts[1]
	setAttrs(ctx, resp, map[string]string{
		"id":         req.ID,
		"repository": repository,
		"uuid":       parts[2],
	})
}

func (r *RepositoryNotificationResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// All meaningful fields use RequiresReplace — update is never called.
}

func (r *RepositoryNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RepositoryNotificationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRepoNotification(ctx, state.Repository.ValueString(), state.UUID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error Deleting Repository Notification", err.Error())
	}
}

func (r *RepositoryNotificationResource) mapToState(ctx context.Context, notif *client.Notification, repository string, state *RepositoryNotificationResourceModel) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	state.Repository = types.StringValue(repository)
	state.UUID = types.StringValue(notif.UUID)
	state.Title = types.StringValue(notif.Title)
	state.Event = types.StringValue(notif.Event)
	state.Method = types.StringValue(notif.Method)
	state.NumberOfFailures = types.Int64Value(int64(notif.NumberOfFailures))
	state.ID = types.StringValue(fmt.Sprintf("%s/%s", repository, notif.UUID))

	configMap, d := anyMapToTFStringMap(ctx, notif.Config)
	diagnostics.Append(d...)
	state.Config = configMap

	eventConfigMap, d := anyMapToTFStringMap(ctx, notif.EventConfig)
	diagnostics.Append(d...)
	state.EventConfig = eventConfigMap

	return diagnostics
}

// mapToStringAny converts a types.Map (string values) to map[string]any.
func mapToStringAny(ctx context.Context, m types.Map) (map[string]any, error) {
	if m.IsNull() || m.IsUnknown() {
		return map[string]any{}, nil
	}
	elems := make(map[string]string)
	if d := m.ElementsAs(ctx, &elems, false); d.HasError() {
		return nil, fmt.Errorf("converting map elements")
	}
	result := make(map[string]any, len(elems))
	for k, v := range elems {
		result[k] = v
	}
	return result, nil
}

// anyMapToTFStringMap converts a map[string]any (with string values) to a types.Map.
func anyMapToTFStringMap(ctx context.Context, m map[string]any) (types.Map, diag.Diagnostics) {
	strMap := make(map[string]string, len(m))
	for k, v := range m {
		strMap[k] = fmt.Sprintf("%v", v)
	}
	return types.MapValueFrom(ctx, types.StringType, strMap)
}
