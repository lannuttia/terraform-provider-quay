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

var _ resource.Resource = &RepositoryUserPermissionResource{}
var _ resource.ResourceWithImportState = &RepositoryUserPermissionResource{}

type RepositoryUserPermissionResource struct {
	client *client.Client
}

type RepositoryUserPermissionResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Repository types.String `tfsdk:"repository"`
	Username   types.String `tfsdk:"username"`
	Role       types.String `tfsdk:"role"`
}

func NewRepositoryUserPermissionResource() resource.Resource {
	return &RepositoryUserPermissionResource{}
}

func (r *RepositoryUserPermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_repository_user_permission"
}

func (r *RepositoryUserPermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a user's permission on a Quay repository.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier (repository/username).",
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
			"username": schema.StringAttribute{
				Required:    true,
				Description: "The username to grant permission to.",
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

func (r *RepositoryUserPermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RepositoryUserPermissionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID: namespace/repo_name/username
	parts, diags := splitImportID(resp, req.ID, 3, "my-org/my-repo/alice")
	if diags.HasError() {
		return
	}
	repository := parts[0] + "/" + parts[1]
	setAttrs(ctx, resp, map[string]string{
		"id":         repository + "/" + parts[2],
		"repository": repository,
		"username":   parts[2],
	})
}

func (r *RepositoryUserPermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RepositoryUserPermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	perm, err := r.client.SetRepoUserPermission(ctx, plan.Repository.ValueString(), plan.Username.ValueString(), client.SetPermissionRequest{
		Role: plan.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Setting Repository User Permission", err.Error())
		return
	}

	plan.Role = types.StringValue(perm.Role)
	plan.ID = types.StringValue(fmt.Sprintf("%s/%s", plan.Repository.ValueString(), plan.Username.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RepositoryUserPermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RepositoryUserPermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	perm, err := r.client.GetRepoUserPermission(ctx, state.Repository.ValueString(), state.Username.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Repository User Permission", err.Error())
		return
	}

	state.Role = types.StringValue(perm.Role)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RepositoryUserPermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RepositoryUserPermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	perm, err := r.client.SetRepoUserPermission(ctx, plan.Repository.ValueString(), plan.Username.ValueString(), client.SetPermissionRequest{
		Role: plan.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Repository User Permission", err.Error())
		return
	}

	plan.Role = types.StringValue(perm.Role)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *RepositoryUserPermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RepositoryUserPermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRepoUserPermission(ctx, state.Repository.ValueString(), state.Username.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error Deleting Repository User Permission", err.Error())
	}
}
