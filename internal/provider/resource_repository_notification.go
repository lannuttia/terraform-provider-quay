package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lannuttia/terraform-provider-quay/internal/client"
)

var _ resource.Resource = &RepositoryNotificationResource{}
var _ resource.ResourceWithImportState = &RepositoryNotificationResource{}
var _ resource.ResourceWithConfigValidators = &RepositoryNotificationResource{}

type RepositoryNotificationResource struct {
	client *client.Client
}

// -------------------------------------------------------------------------
// Terraform state model
// -------------------------------------------------------------------------

type RepositoryNotificationResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Repository       types.String `tfsdk:"repository"`
	UUID             types.String `tfsdk:"uuid"`
	Title            types.String `tfsdk:"title"`
	Event            types.String `tfsdk:"event"`
	NumberOfFailures types.Int64  `tfsdk:"number_of_failures"`

	// Exactly one method block must be set.
	Webhook          *WebhookConfig          `tfsdk:"webhook"`
	Slack            *SlackConfig            `tfsdk:"slack"`
	Email            *EmailConfig            `tfsdk:"email"`
	Hipchat          *HipchatConfig          `tfsdk:"hipchat"`
	Flowdock         *FlowdockConfig         `tfsdk:"flowdock"`
	QuayNotification *QuayNotificationConfig `tfsdk:"quay_notification"`

	// At most one event config block may be set, depending on the event.
	VulnerabilityConfig *VulnerabilityEventConfig `tfsdk:"vulnerability_config"`
	BuildConfig         *BuildEventConfig         `tfsdk:"build_config"`
}

// Method config types

type WebhookConfig struct {
	URL types.String `tfsdk:"url"`
}

type SlackConfig struct {
	URL types.String `tfsdk:"url"`
}

type EmailConfig struct {
	Email types.String `tfsdk:"email"`
}

type HipchatConfig struct {
	NotificationToken types.String `tfsdk:"notification_token"`
	RoomID            types.String `tfsdk:"room_id"`
}

type FlowdockConfig struct {
	FlowAPIToken types.String `tfsdk:"flow_api_token"`
}

type QuayNotificationConfig struct {
	// Quay internal notifications carry no extra config fields.
}

// Event config types

type VulnerabilityEventConfig struct {
	Level types.String `tfsdk:"level"`
}

type BuildEventConfig struct {
	RefRegex types.String `tfsdk:"ref_regex"`
}

// -------------------------------------------------------------------------
// Constructor / metadata
// -------------------------------------------------------------------------

func NewRepositoryNotificationResource() resource.Resource {
	return &RepositoryNotificationResource{}
}

func (r *RepositoryNotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_notification"
}

// -------------------------------------------------------------------------
// Schema
// -------------------------------------------------------------------------

func (r *RepositoryNotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	requiresReplace := []planmodifier.String{stringplanmodifier.RequiresReplace()}

	resp.Schema = schema.Schema{
		Description: "Manages a notification on a Quay repository. " +
			"Because the Quay API has no update endpoint for notifications, any change to notification " +
			"attributes will cause the existing notification to be deleted and a new one created. " +
			"Exactly one method block (webhook, slack, email, hipchat, flowdock, or quay_notification) " +
			"must be provided. For vulnerability_found events supply a vulnerability_config block; " +
			"for build events an optional build_config block may be supplied.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier (repository/uuid).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"repository": schema.StringAttribute{
				Required:      true,
				Description:   "The repository this notification belongs to, in namespace/name format.",
				PlanModifiers: requiresReplace,
			},
			"uuid": schema.StringAttribute{
				Computed:    true,
				Description: "The server-assigned UUID of the notification.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Description:   "A human-readable title for the notification.",
				PlanModifiers: append(requiresReplace, stringplanmodifier.UseStateForUnknown()),
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
				PlanModifiers: requiresReplace,
			},
			"number_of_failures": schema.Int64Attribute{
				Computed:    true,
				Description: "The number of consecutive delivery failures.",
			},

			// ----- method blocks -----

			"webhook": schema.SingleNestedAttribute{
				Optional:      true,
				Description:   "Deliver notifications via an HTTP POST to a URL.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Required:      true,
						Description:   "The URL to POST the notification payload to.",
						PlanModifiers: requiresReplace,
					},
				},
			},
			"slack": schema.SingleNestedAttribute{
				Optional:      true,
				Description:   "Deliver notifications to a Slack channel via an Incoming Webhook URL.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Required:      true,
						Description:   "The Slack Incoming Webhook URL.",
						PlanModifiers: requiresReplace,
					},
				},
			},
			"email": schema.SingleNestedAttribute{
				Optional:      true,
				Description:   "Deliver notifications by email.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
				Attributes: map[string]schema.Attribute{
					"email": schema.StringAttribute{
						Required:      true,
						Description:   "The recipient email address.",
						PlanModifiers: requiresReplace,
					},
				},
			},
			"hipchat": schema.SingleNestedAttribute{
				Optional:      true,
				Description:   "Deliver notifications to a HipChat room.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
				Attributes: map[string]schema.Attribute{
					"notification_token": schema.StringAttribute{
						Required:      true,
						Description:   "The HipChat API notification token.",
						PlanModifiers: requiresReplace,
					},
					"room_id": schema.StringAttribute{
						Required:      true,
						Description:   "The HipChat room ID.",
						PlanModifiers: requiresReplace,
					},
				},
			},
			"flowdock": schema.SingleNestedAttribute{
				Optional:      true,
				Description:   "Deliver notifications to a Flowdock flow.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
				Attributes: map[string]schema.Attribute{
					"flow_api_token": schema.StringAttribute{
						Required:      true,
						Description:   "The Flowdock flow API token.",
						PlanModifiers: requiresReplace,
					},
				},
			},
			"quay_notification": schema.SingleNestedAttribute{
				Optional:      true,
				Description:   "Deliver notifications via the built-in Quay notification system.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
				Attributes:    map[string]schema.Attribute{},
			},

			// ----- event config blocks -----

			"vulnerability_config": schema.SingleNestedAttribute{
				Optional:      true,
				Description:   "Event configuration for vulnerability_found events.",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
				Attributes: map[string]schema.Attribute{
					"level": schema.StringAttribute{
						Required:    true,
						Description: "The minimum vulnerability severity level that triggers the notification. One of: Low, Medium, High, Critical.",
						Validators: []validator.String{
							stringvalidator.OneOf("Low", "Medium", "High", "Critical"),
						},
						PlanModifiers: requiresReplace,
					},
				},
			},
			"build_config": schema.SingleNestedAttribute{
				Optional:      true,
				Description:   "Event configuration for build events (build_queued, build_start, build_success, build_failure, build_cancelled).",
				PlanModifiers: []planmodifier.Object{objectplanmodifier.RequiresReplace()},
				Attributes: map[string]schema.Attribute{
					"ref_regex": schema.StringAttribute{
						Optional:      true,
						Description:   "A regular expression that Git refs must match to trigger the notification.",
						PlanModifiers: requiresReplace,
					},
				},
			},
		},
	}
}

// ConfigValidators enforces that exactly one method block is set.
func (r *RepositoryNotificationResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("webhook"),
			path.MatchRoot("slack"),
			path.MatchRoot("email"),
			path.MatchRoot("hipchat"),
			path.MatchRoot("flowdock"),
			path.MatchRoot("quay_notification"),
		),
	}
}

// -------------------------------------------------------------------------
// CRUD
// -------------------------------------------------------------------------

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

	method, config := plan.methodAndConfig()
	eventConfig := plan.eventConfig()

	notif, err := r.client.CreateRepoNotification(ctx, plan.Repository.ValueString(), client.CreateNotificationRequest{
		Event:       plan.Event.ValueString(),
		Method:      method,
		Config:      config,
		EventConfig: eventConfig,
		Title:       plan.Title.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Repository Notification", err.Error())
		return
	}

	resp.Diagnostics.Append(r.mapToState(notif, plan.Repository.ValueString(), &plan)...)
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

	resp.Diagnostics.Append(r.mapToState(notif, state.Repository.ValueString(), &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
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

func (r *RepositoryNotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID: namespace/repo_name/uuid  (3 parts, last part is the UUID)
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

// -------------------------------------------------------------------------
// Mapping helpers
// -------------------------------------------------------------------------

// methodAndConfig derives the API method string and config map from whichever
// method block is set in the plan.
func (m *RepositoryNotificationResourceModel) methodAndConfig() (string, map[string]any) {
	switch {
	case m.Webhook != nil:
		return "webhook", map[string]any{"url": m.Webhook.URL.ValueString()}
	case m.Slack != nil:
		return "slack", map[string]any{"url": m.Slack.URL.ValueString()}
	case m.Email != nil:
		return "email", map[string]any{"email": m.Email.Email.ValueString()}
	case m.Hipchat != nil:
		return "hipchat", map[string]any{
			"notification_token": m.Hipchat.NotificationToken.ValueString(),
			"room_id":            m.Hipchat.RoomID.ValueString(),
		}
	case m.Flowdock != nil:
		return "flowdock", map[string]any{"flow_api_token": m.Flowdock.FlowAPIToken.ValueString()}
	case m.QuayNotification != nil:
		return "quay_notification", map[string]any{}
	default:
		// ConfigValidators guarantees this is unreachable.
		return "", map[string]any{}
	}
}

// eventConfig derives the API event_config map from whichever event config
// block is set in the plan (or an empty map if none is set).
func (m *RepositoryNotificationResourceModel) eventConfig() map[string]any {
	switch {
	case m.VulnerabilityConfig != nil:
		return map[string]any{"level": m.VulnerabilityConfig.Level.ValueString()}
	case m.BuildConfig != nil:
		cfg := map[string]any{}
		if !m.BuildConfig.RefRegex.IsNull() && !m.BuildConfig.RefRegex.IsUnknown() {
			cfg["ref-regex"] = m.BuildConfig.RefRegex.ValueString()
		}
		return cfg
	default:
		return map[string]any{}
	}
}

// mapToState populates the model from a Notification API response.
func (r *RepositoryNotificationResource) mapToState(notif *client.Notification, repository string, state *RepositoryNotificationResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	state.Repository = types.StringValue(repository)
	state.UUID = types.StringValue(notif.UUID)
	state.Title = types.StringValue(notif.Title)
	state.Event = types.StringValue(notif.Event)
	state.NumberOfFailures = types.Int64Value(int64(notif.NumberOfFailures))
	state.ID = types.StringValue(fmt.Sprintf("%s/%s", repository, notif.UUID))

	// Clear all method blocks, then populate only the one that matches.
	state.Webhook = nil
	state.Slack = nil
	state.Email = nil
	state.Hipchat = nil
	state.Flowdock = nil
	state.QuayNotification = nil

	switch notif.Method {
	case "webhook":
		state.Webhook = &WebhookConfig{URL: types.StringValue(stringFromConfig(notif.Config, "url"))}
	case "slack":
		state.Slack = &SlackConfig{URL: types.StringValue(stringFromConfig(notif.Config, "url"))}
	case "email":
		state.Email = &EmailConfig{Email: types.StringValue(stringFromConfig(notif.Config, "email"))}
	case "hipchat":
		state.Hipchat = &HipchatConfig{
			NotificationToken: types.StringValue(stringFromConfig(notif.Config, "notification_token")),
			RoomID:            types.StringValue(stringFromConfig(notif.Config, "room_id")),
		}
	case "flowdock":
		state.Flowdock = &FlowdockConfig{FlowAPIToken: types.StringValue(stringFromConfig(notif.Config, "flow_api_token"))}
	case "quay_notification":
		state.QuayNotification = &QuayNotificationConfig{}
	default:
		diags.AddError("Unknown Notification Method",
			fmt.Sprintf("The Quay API returned an unrecognised notification method %q. "+
				"Please report this as a provider bug.", notif.Method))
		return diags
	}

	// Clear all event config blocks, then populate only the one that matches.
	state.VulnerabilityConfig = nil
	state.BuildConfig = nil

	switch {
	case len(notif.EventConfig) > 0 && notif.EventConfig["level"] != nil:
		state.VulnerabilityConfig = &VulnerabilityEventConfig{
			Level: types.StringValue(fmt.Sprintf("%v", notif.EventConfig["level"])),
		}
	case len(notif.EventConfig) > 0 && notif.EventConfig["ref-regex"] != nil:
		state.BuildConfig = &BuildEventConfig{
			RefRegex: types.StringValue(fmt.Sprintf("%v", notif.EventConfig["ref-regex"])),
		}
	}

	return diags
}

// stringFromConfig safely extracts a string value from a map[string]any.
func stringFromConfig(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}
