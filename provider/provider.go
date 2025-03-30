package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	providerSchema "github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type exampleProvider struct{}

func New() provider.Provider {
	return &exampleProvider{}
}

func (p *exampleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "example"
}

func (p *exampleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = providerSchema.Schema{}
}

func (p *exampleProvider) Configure(_ context.Context, _ provider.ConfigureRequest, _ *provider.ConfigureResponse) {
	// No-op for this example
}

func (p *exampleProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil // No resources implemented in this example
}

type helloWorldDataSource struct{}

func (p *exampleProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return &helloWorldDataSource{} },
	}
}

func (d *helloWorldDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "hello_world"
}

func (d *helloWorldDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			"message": datasourceSchema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *helloWorldDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	resp.Diagnostics.Append(resp.State.Set(ctx, map[string]types.String{
		"message": types.StringValue("Hello, World!"),
	})...)
}
