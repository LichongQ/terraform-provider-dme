package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"terraform-provider-dme/internal/provider"
)

// main is the entry point for the Terraform DME provider.
// It registers the provider implementation with the Terraform plugin SDK.
func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.New,
	})
}
