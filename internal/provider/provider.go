// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure discordProvider satisfies various provider interfaces.
var _ provider.Provider = &discordProvider{}

// discordProvider defines the provider implementation.
type discordProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// discordProviderModel describes the provider data model.
type discordProviderModel struct {
	AuthenticationToken types.String `tfsdk:"authentication_token"`
}

func (p *discordProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "discord"
	resp.Version = p.version
}

func (p *discordProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"authentication_token": schema.StringAttribute{
				MarkdownDescription: "authentication token for the discord bot",
				Optional:            true,
			},
		},
	}
}

func (p *discordProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data discordProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if data.AuthenticationToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("authentication_token"),
			"Unknown authentication token",
			"please provide a valid Discord API token",
		)
	}
	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	AuthenticationToken := os.Getenv("DISCORD_AUTHENTICATION_TOKEN")

	if !data.AuthenticationToken.IsNull() {
		AuthenticationToken = data.AuthenticationToken.ValueString()
	}

	if AuthenticationToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("authentication_token"),
			"Missing authentication_token",
			"Please add a authentication token or set the env var DISCORD_AUTHENTICATION_TOKEN",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}
	client, err := discordgo.New("Bot " + AuthenticationToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Discord API Client",
			"An unexpected error occurred when creating the Discord API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Discord Client Error: "+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *discordProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *discordProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewServerDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &discordProvider{
			version: version,
		}
	}
}
