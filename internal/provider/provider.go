// Package provider implements the Terraform provider for eDME (Enterprise Data Management Engine).
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"terraform-provider-dme/internal/client"
	"terraform-provider-dme/internal/datasource"
	"terraform-provider-dme/internal/resource"
)

// New returns a new Terraform provider instance for eDME.
func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "eDME management endpoint URL, e.g. `https://10.0.0.1:26335`",
			},
			"user_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "eDME northbound user name",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "eDME northbound user password",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dme_vm": resource.ResourceVM(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dme_vm":  datasource.DataSourceVM(),
			"dme_vms": datasource.DataSourceVMs(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

// configureProvider initializes the eDME API client and authenticates.
func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	cfg := client.Config{
		Endpoint: d.Get("endpoint").(string),
		UserName: d.Get("user_name").(string),
		Password: d.Get("password").(string),
	}

	cl := client.NewClient(cfg)

	// Authenticate to obtain the initial session token.
	if err := cl.Authenticate(ctx); err != nil {
		return nil, diag.Errorf("failed to authenticate with eDME: %s", err)
	}

	return cl, nil
}
