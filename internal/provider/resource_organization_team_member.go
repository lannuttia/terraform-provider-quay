package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lannuttia/terraform-provider-quay/internal/client"
)

var _ resource.Resource = &OrganizationTeamMemberResource{}

type OrganizationTeamMemberResource struct {
	client *client.Client
}

type OrganizationTeamMemberResourceModel struct {
	ID         types.String `tfsdk:"id"`
	OrgName    types.String `tfsdk:"org_name"`
	TeamName   types.String `tfsdk:"team_name"`
	MemberName types.String `tfsdk:"member_name"`
	IsRobot    types.Bool   `tfsdk:"is_robot"`
}

func NewOrganizationTeamMemberResource() resource.Resource {
	return &OrganizationTeamMemberResource{}
}

func (r *OrganizationTeamMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_team_member"
}

func (r *OrganizationTeamMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages membership of a user or robot in a Quay organization team.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier (org_name/team_name/member_name).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization.",
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
			"member_name": schema.StringAttribute{
				Required:    true,
				Description: "The username or robot account name to add to the team.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_robot": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the member is a robot account.",
			},
		},
	}
}

func (r *OrganizationTeamMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationTeamMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationTeamMemberResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AddOrgTeamMember(ctx, plan.OrgName.ValueString(), plan.TeamName.ValueString(), plan.MemberName.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error Adding Team Member", err.Error())
		return
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s/%s/%s", plan.OrgName.ValueString(), plan.TeamName.ValueString(), plan.MemberName.ValueString()))
	plan.IsRobot = types.BoolValue(strings.Contains(plan.MemberName.ValueString(), "+"))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationTeamMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationTeamMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	members, err := r.client.GetOrgTeamMembers(ctx, state.OrgName.ValueString(), state.TeamName.ValueString(), true)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Team Members", err.Error())
		return
	}

	for _, m := range members {
		if m.Name == state.MemberName.ValueString() {
			state.IsRobot = types.BoolValue(m.IsRobot)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// Member not found — remove from state
	resp.State.RemoveResource(ctx)
}

func (r *OrganizationTeamMemberResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
	// All fields use RequiresReplace — update is never called.
}

func (r *OrganizationTeamMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationTeamMemberResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.RemoveOrgTeamMember(ctx, state.OrgName.ValueString(), state.TeamName.ValueString(), state.MemberName.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error Removing Team Member", err.Error())
	}
}
