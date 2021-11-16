package pulumi

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("PULUMI_ACCESS_TOKEN", nil),
				Description: "The access token for Pulumi's cloud",
			},
		},
		ResourcesMap: map[string]*schema.Resource{},
		DataSourcesMap: map[string]*schema.Resource{
			"pulumi_stack_outputs": dataSourceStackOutputs(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	pulumi_token := d.Get("token").(string)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return pulumi_token, diags
}
