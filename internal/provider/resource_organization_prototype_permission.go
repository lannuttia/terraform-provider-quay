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

var _ resource.Resource = &OrganizationPrototypePermissionResource{}

type OrganizationPrototypePermissionResource struct {
	client *client.Client
}

type OrganizationPrototypePermissionResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	OrgName              types.String `tfsdk:"org_name"`
	Role                 types.String `tfsdk:"role"`
	DelegateKind         types.String `tfsdk:"delegate_kind"`
	DelegateName         types.String `tfsdk:"delegate_name"`
	ActivatingUserName   types.String `tfsdk:"activating_user_name"`
}

func NewOrganizationPrototypePermissionResource() resource.Resource {
	return &OrganizationPrototypePermissionResource{}
}

func (r *OrganizationPrototypePermissionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_prototype_permission"
}

func (r *OrganizationPrototypePermissionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a default permission prototype for a Quay organization. Prototypes automatically grant permissions when new repositories are created.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The server-assigned prototype ID.",
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
			"role": schema.StringAttribute{
				Required:    true,
				Description: "The role to grant. One of: read, write, admin.",
				Validators: []validator.String{
					stringvalidator.OneOf("read", "write", "admin"),
				},
			},
			"delegate_kind": schema.StringAttribute{
				Required:    true,
				Description: "Whether the delegate is a 'user' or 'team'.",
				Validators: []validator.String{
					stringvalidator.OneOf("user", "team"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"delegate_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the user or team to grant the permission to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"activating_user_name": schema.StringAttribute{
				Optional:    true,
				Description: "If set, this prototype only applies to repositories created by this specific user.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

func (r *OrganizationPrototypePermissionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationPrototypePermissionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationPrototypePermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreatePrototypeRequest{
		Role: plan.Role.ValueString(),
		Delegate: client.PrototypeDelegate{
			Kind: plan.DelegateKind.ValueString(),
			Name: plan.DelegateName.ValueString(),
		},
	}
	if !plan.ActivatingUserName.IsNull() && plan.ActivatingUserName.ValueString() != "" {
		createReq.ActivatingUser = &client.PrototypeActivatingUser{Name: plan.ActivatingUserName.ValueString()}
	}

	proto, err := r.client.CreateOrgPrototype(ctx, plan.OrgName.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Prototype Permission", err.Error())
		return
	}

	r.mapToState(proto, plan.OrgName.ValueString(), &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationPrototypePermissionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationPrototypePermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	proto, err := r.client.GetOrgPrototype(ctx, state.OrgName.ValueString(), state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Prototype Permission", err.Error())
		return
	}

	r.mapToState(proto, state.OrgName.ValueString(), &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationPrototypePermissionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationPrototypePermissionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	proto, err := r.client.UpdateOrgPrototype(ctx, plan.OrgName.ValueString(), plan.ID.ValueString(), client.UpdatePrototypeRequest{
		Role: plan.Role.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Prototype Permission", err.Error())
		return
	}

	r.mapToState(proto, plan.OrgName.ValueString(), &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationPrototypePermissionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationPrototypePermissionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrgPrototype(ctx, state.OrgName.ValueString(), state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error Deleting Prototype Permission", err.Error())
	}
}

func (r *OrganizationPrototypePermissionResource) mapToState(proto *client.Prototype, orgName string, state *OrganizationPrototypePermissionResourceModel) {
	state.ID = types.StringValue(proto.ID)
	state.OrgName = types.StringValue(orgName)
	state.Role = types.StringValue(proto.Role)
	state.DelegateKind = types.StringValue(proto.Delegate.Kind)
	state.DelegateName = types.StringValue(proto.Delegate.Name)
	if proto.ActivatingUser != nil {
		state.ActivatingUserName = types.StringValue(proto.ActivatingUser.Name)
	} else {
		state.ActivatingUserName = types.StringNull()
	}
}
