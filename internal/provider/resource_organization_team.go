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

var _ resource.Resource = &OrganizationTeamResource{}

type OrganizationTeamResource struct {
	client *client.Client
}

type OrganizationTeamResourceModel struct {
	ID          types.String `tfsdk:"id"`
	OrgName     types.String `tfsdk:"org_name"`
	Name        types.String `tfsdk:"name"`
	Role        types.String `tfsdk:"role"`
	Description types.String `tfsdk:"description"`
}

func NewOrganizationTeamResource() resource.Resource {
	return &OrganizationTeamResource{}
}

func (r *OrganizationTeamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_team"
}

func (r *OrganizationTeamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a team within a Quay organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier (org_name/team_name).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization that owns this team.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the team.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"role": schema.StringAttribute{
				Required:    true,
				Description: "The org-wide role for this team. One of: member, creator, admin.",
				Validators: []validator.String{
					stringvalidator.OneOf("member", "creator", "admin"),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A markdown description of the team.",
			},
		},
	}
}

func (r *OrganizationTeamResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationTeamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationTeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := r.client.UpsertOrgTeam(ctx, plan.OrgName.ValueString(), plan.Name.ValueString(), client.UpsertTeamRequest{
		Role:        plan.Role.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Organization Team", err.Error())
		return
	}

	r.mapToState(team, plan.OrgName.ValueString(), plan.Name.ValueString(), &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationTeamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationTeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The Quay API has no GET single team endpoint — read via member list and org data.
	// We'll confirm existence by fetching members (returns 404 if team doesn't exist).
	_, err := r.client.GetOrgTeamMembers(ctx, state.OrgName.ValueString(), state.Name.ValueString(), false)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Organization Team", err.Error())
		return
	}

	// State remains as-is since we can't read role/description back from the API
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationTeamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationTeamResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	team, err := r.client.UpsertOrgTeam(ctx, plan.OrgName.ValueString(), plan.Name.ValueString(), client.UpsertTeamRequest{
		Role:        plan.Role.ValueString(),
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Organization Team", err.Error())
		return
	}

	r.mapToState(team, plan.OrgName.ValueString(), plan.Name.ValueString(), &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationTeamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationTeamResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrgTeam(ctx, state.OrgName.ValueString(), state.Name.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error Deleting Organization Team", err.Error())
	}
}

func (r *OrganizationTeamResource) mapToState(team *client.Team, orgName, teamName string, state *OrganizationTeamResourceModel) {
	state.OrgName = types.StringValue(orgName)
	state.Name = types.StringValue(teamName)
	if team.Role != "" {
		state.Role = types.StringValue(team.Role)
	}
	if team.Description != "" {
		state.Description = types.StringValue(team.Description)
	} else if state.Description.IsNull() {
		state.Description = types.StringValue("")
	}
	state.ID = types.StringValue(fmt.Sprintf("%s/%s", orgName, teamName))
}
