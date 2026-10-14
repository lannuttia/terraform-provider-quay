package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lannuttia/terraform-provider-quay/internal/client"
)

var _ resource.Resource = &RepositoryTeamPermissionResource{}
var _ resource.ResourceWithImportState = &RepositoryTeamPermissionResource{}

type RepositoryTeamPermissionResource struct {
	client *client.Client
}

type RepositoryTeamPermissionResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Repository types.String `tfsdk:"repository"`
	TeamName   types.String `tfsdk:"team_name"`
	Role       types.String `tfsdk:"role"`
}

func NewRepositoryTeamPermissionResource() resource.Resource {
	return &RepositoryTeamPermissionResource{}
}

func (r *RepositoryTeamPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_team_permission"
}

func (r *RepositoryTeamPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a team's permission on a Quay repository.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier (repository/team_name).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"repository": schema.StringAttribute{
				Required:    true,
				Description: "The repository in namespace/name format.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"team_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the team.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Required:    true,
				Description: "The permission role. One of: read, write, admin.",
				Validators: []validator.String{
					stringvalidator.OneOf("read", "write", "admin"),
				},
			},
		},
	}
}

func (r *RepositoryTeamPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RepositoryTeamPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID: namespace/repo_name/team_name
	parts, diags := splitImportID(resp, req.ID, 3, "my-org/my-repo/my-team")
	if diags.HasError() {
		return
	}
	repository := parts[0] + "/" + parts[1]
	setAttrs(ctx, resp, map[string]string{
		"id":         repository + "/" + parts[2],
		"repository": repository,
		"team_name":  parts[2],
	})
}

func (r *RepositoryTeamPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RepositoryTeamPermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	perm, err := r.client.SetRepoTeamPermission(ctx, plan.Repository.ValueString(), plan.TeamName.ValueString(), client.SetPermissionRequest{
		Role: plan.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Setting Repository Team Permission", err.Error())
		return
	}

	plan.Role = types.StringValue(perm.Role)
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.Repository.ValueString(), plan.TeamName.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RepositoryTeamPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RepositoryTeamPermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	perm, err := r.client.GetRepoTeamPermission(ctx, state.Repository.ValueString(), state.TeamName.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Repository Team Permission", err.Error())
		return
	}

	state.Role = types.StringValue(perm.Role)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RepositoryTeamPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RepositoryTeamPermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	perm, err := r.client.SetRepoTeamPermission(ctx, plan.Repository.ValueString(), plan.TeamName.ValueString(), client.SetPermissionRequest{
		Role: plan.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Repository Team Permission", err.Error())
		return
	}

	plan.Role = types.StringValue(perm.Role)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RepositoryTeamPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RepositoryTeamPermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRepoTeamPermission(ctx, state.Repository.ValueString(), state.TeamName.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error Deleting Repository Team Permission", err.Error())
	}
}
