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

var _ resource.Resource = &OrganizationRobotResource{}
var _ resource.ResourceWithImportState = &OrganizationRobotResource{}

type OrganizationRobotResource struct {
	client *client.Client
}

type OrganizationRobotResourceModel struct {
	ID          types.String `tfsdk:"id"`
	OrgName     types.String `tfsdk:"org_name"`
	ShortName   types.String `tfsdk:"short_name"`
	FullName    types.String `tfsdk:"full_name"`
	Token       types.String `tfsdk:"token"`
	Description types.String `tfsdk:"description"`
}

func NewOrganizationRobotResource() resource.Resource {
	return &OrganizationRobotResource{}
}

func (r *OrganizationRobotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_robot"
}

func (r *OrganizationRobotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a robot account within a Quay organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier (org_name/short_name).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization that owns this robot.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"short_name": schema.StringAttribute{
				Required:    true,
				Description: "The short name of the robot (without the org prefix).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"full_name": schema.StringAttribute{
				Computed:    true,
				Description: "The full robot account name (org_name+short_name).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"token": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The robot account's credentials token. Only available immediately after creation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the robot account.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *OrganizationRobotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationRobotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationRobotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	robot, err := r.client.CreateOrgRobot(ctx, plan.OrgName.ValueString(), plan.ShortName.ValueString(), client.CreateRobotRequest{
		Description: plan.Description.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Organization Robot", err.Error())
		return
	}

	r.mapToState(robot, plan.OrgName.ValueString(), &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationRobotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationRobotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	robot, err := r.client.GetOrgRobot(ctx, state.OrgName.ValueString(), state.ShortName.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Organization Robot", err.Error())
		return
	}

	// Preserve token from state — the API never returns it after creation
	robot.Token = state.Token.ValueString()
	r.mapToState(robot, state.OrgName.ValueString(), &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationRobotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Robot accounts have no update API — only short_name/org_name trigger replacement.
	// Description updates require delete+recreate (RequiresReplace is not set so we can
	// at least keep state consistent; in practice only non-replacing fields would reach here).
	var plan OrganizationRobotResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationRobotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationRobotResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrgRobot(ctx, state.OrgName.ValueString(), state.ShortName.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error Deleting Organization Robot", err.Error())
	}
}

func (r *OrganizationRobotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID: org_name/short_name
	parts, diags := splitImportID(resp, req.ID, 2, "my-org/my-robot")
	if diags.HasError() {
		return
	}
	setAttrs(ctx, resp, map[string]string{
		"id":         req.ID,
		"org_name":   parts[0],
		"short_name": parts[1],
	})
}

func (r *OrganizationRobotResource) mapToState(robot *client.Robot, orgName string, state *OrganizationRobotResourceModel) {
	state.OrgName = types.StringValue(orgName)
	// Derive short name from full name if available
	shortName := state.ShortName.ValueString()
	if robot.Name != "" {
		parts := strings.SplitN(robot.Name, "+", 2)
		if len(parts) == 2 {
			shortName = parts[1]
		}
	}
	state.ShortName = types.StringValue(shortName)
	state.FullName = types.StringValue(robot.Name)
	state.Description = types.StringValue(robot.Description)
	if robot.Token != "" {
		state.Token = types.StringValue(robot.Token)
	}
	state.ID = types.StringValue(fmt.Sprintf("%s/%s", orgName, shortName))
}
