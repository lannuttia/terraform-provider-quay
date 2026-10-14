package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lannuttia/terraform-provider-quay/internal/client"
)

var _ datasource.DataSource = &UserDataSource{}

type UserDataSource struct {
	client *client.Client
}

type UserDataSourceModel struct {
	Username  types.String `tfsdk:"username"`
	Email     types.String `tfsdk:"email"`
	Verified  types.Bool   `tfsdk:"verified"`
	Anonymous types.Bool   `tfsdk:"anonymous"`
}

func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetches information about a Quay user. If username is omitted or empty, returns the currently authenticated user.",
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The username to look up. If omitted, returns the currently authenticated user.",
			},
			"email": schema.StringAttribute{
				Computed:    true,
				Description: "The user's email address.",
			},
			"verified": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the user's email is verified.",
			},
			"anonymous": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether this is an anonymous (unauthenticated) user.",
			},
		},
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *client.Client, got %T", req.ProviderData))
		return
	}
	d.client = c
}

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state UserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var user *client.User
	var err error

	if state.Username.IsNull() || state.Username.IsUnknown() || state.Username.ValueString() == "" {
		user, err = d.client.GetCurrentUser(ctx)
	} else {
		user, err = d.client.GetUser(ctx, state.Username.ValueString())
	}

	if err != nil {
		resp.Diagnostics.AddError("Error Reading User", err.Error())
		return
	}

	state.Username = types.StringValue(user.Username)
	state.Email = types.StringValue(user.Email)
	state.Verified = types.BoolValue(user.Verified)
	state.Anonymous = types.BoolValue(user.Anonymous)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
