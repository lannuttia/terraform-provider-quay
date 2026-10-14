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

var _ resource.Resource = &OrganizationApplicationResource{}
var _ resource.ResourceWithImportState = &OrganizationApplicationResource{}

type OrganizationApplicationResource struct {
	client *client.Client
}

type OrganizationApplicationResourceModel struct {
	ID             types.String `tfsdk:"id"`
	OrgName        types.String `tfsdk:"org_name"`
	ClientID       types.String `tfsdk:"client_id"`
	ClientSecret   types.String `tfsdk:"client_secret"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	ApplicationURI types.String `tfsdk:"application_uri"`
	RedirectURI    types.String `tfsdk:"redirect_uri"`
	AvatarEmail    types.String `tfsdk:"avatar_email"`
}

func NewOrganizationApplicationResource() resource.Resource {
	return &OrganizationApplicationResource{}
}

func (r *OrganizationApplicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organization_application"
}

func (r *OrganizationApplicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an OAuth application within a Quay organization.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource identifier (org_name/client_id).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"org_name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the organization that owns this application.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"client_id": schema.StringAttribute{
				Computed:    true,
				Description: "The OAuth client ID assigned by Quay.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_secret": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The OAuth client secret assigned by Quay.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The application name.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "A description of the application.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"application_uri": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The URI of the application's homepage.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"redirect_uri": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The OAuth redirect URI.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"avatar_email": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The email address used for the application's avatar.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *OrganizationApplicationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationApplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.CreateOrgApplication(ctx, plan.OrgName.ValueString(), client.CreateOrgApplicationRequest{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		ApplicationURI: plan.ApplicationURI.ValueString(),
		RedirectURI:    plan.RedirectURI.ValueString(),
		AvatarEmail:    plan.AvatarEmail.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Organization Application", err.Error())
		return
	}

	r.mapToState(app, plan.OrgName.ValueString(), &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	orgName, clientID := parseOrgAppID(state.ID.ValueString())

	app, err := r.client.GetOrgApplication(ctx, orgName, clientID)
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Organization Application", err.Error())
		return
	}

	r.mapToState(app, orgName, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *OrganizationApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationApplicationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.client.UpdateOrgApplication(ctx, plan.OrgName.ValueString(), plan.ClientID.ValueString(), client.UpdateOrgApplicationRequest{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		ApplicationURI: plan.ApplicationURI.ValueString(),
		RedirectURI:    plan.RedirectURI.ValueString(),
		AvatarEmail:    plan.AvatarEmail.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Error Updating Organization Application", err.Error())
		return
	}

	r.mapToState(app, plan.OrgName.ValueString(), &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OrganizationApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationApplicationResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteOrgApplication(ctx, state.OrgName.ValueString(), state.ClientID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error Deleting Organization Application", err.Error())
	}
}

func (r *OrganizationApplicationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import ID: org_name/client_id
	parts, diags := splitImportID(resp, req.ID, 2, "my-org/abc123clientid")
	if diags.HasError() {
		return
	}
	setAttrs(ctx, resp, map[string]string{
		"id":        req.ID,
		"org_name":  parts[0],
		"client_id": parts[1],
	})
}

func (r *OrganizationApplicationResource) mapToState(app *client.OrgApplication, orgName string, state *OrganizationApplicationResourceModel) {
	state.OrgName = types.StringValue(orgName)
	state.ClientID = types.StringValue(app.ClientID)
	state.ClientSecret = types.StringValue(app.ClientSecret)
	state.Name = types.StringValue(app.Name)
	state.Description = types.StringValue(app.Description)
	state.ApplicationURI = types.StringValue(app.ApplicationURI)
	state.RedirectURI = types.StringValue(app.RedirectURI)
	state.AvatarEmail = types.StringValue(app.AvatarEmail)
	state.ID = types.StringValue(fmt.Sprintf("%s/%s", orgName, app.ClientID))
}

func parseOrgAppID(id string) (orgName, clientID string) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return id, ""
}
