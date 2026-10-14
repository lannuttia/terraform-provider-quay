package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lannuttia/terraform-provider-quay/internal/client"
)

var _ provider.Provider = &QuayProvider{}

type QuayProvider struct {
	version string
}

type QuayProviderModel struct {
	URL   types.String `tfsdk:"url"`
	Token types.String `tfsdk:"token"`
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &QuayProvider{version: version}
	}
}

func (p *QuayProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "quay"
	resp.Version = p.version
}

func (p *QuayProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Quay.io or a self-hosted Red Hat Quay instance.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Optional:    true,
				Description: "The base URL of the Quay instance. Defaults to https://quay.io. Can also be set via the QUAY_URL environment variable.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "OAuth bearer token for authenticating with the Quay API. Can also be set via the QUAY_TOKEN environment variable.",
			},
		},
	}
}

func (p *QuayProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config QuayProviderModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := os.Getenv("QUAY_URL")
	if !config.URL.IsNull() && !config.URL.IsUnknown() {
		url = config.URL.ValueString()
	}

	token := os.Getenv("QUAY_TOKEN")
	if !config.Token.IsNull() && !config.Token.IsUnknown() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing Quay API Token",
			"A Quay API token must be provided via the `token` provider attribute or the QUAY_TOKEN environment variable.",
		)
		return
	}

	c := client.New(url, token)
	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *QuayProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
		NewOrganizationApplicationResource,
		NewOrganizationRobotResource,
		NewOrganizationTeamResource,
		NewOrganizationTeamMemberResource,
		NewOrganizationPrototypePermissionResource,
		NewRepositoryResource,
		NewRepositoryNotificationResource,
		NewRepositoryTeamPermissionResource,
		NewRepositoryUserPermissionResource,
	}
}

func (p *QuayProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationDataSource,
		NewRepositoryDataSource,
		NewUserDataSource,
	}
}

func (p *QuayProvider) Functions(_ context.Context) []func() function.Function {
	return nil
}
