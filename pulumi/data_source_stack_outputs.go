package pulumi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Resource struct {
	urn     string
	outputs map[string]string
}

func dataSourceStackOutputsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := &http.Client{Timeout: 10 * time.Second}

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// The ID will be the org, project, and stack names separated with slashes
	project := d.Get("project").(string)
	stack := d.Get("stack").(string)
	id := fmt.Sprintf("%s/%s/%s",
		d.Get("organization").(string),
		project,
		stack,
	)
	d.SetId(id)

	// Construct the HTTP Request
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.pulumi.com/api/stacks/%s/export", id), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Add pulumi auth
	auth := "token " + m.(string)
	req.Header.Add("Authorization", auth)

	// Send the request
	res, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	// body, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	return diag.FromErr(err)
	// }

	// Parse the response
	response := make(map[string]interface{}, 0)
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return diag.FromErr(err)
	}

	// Find the stack outputs
	stack_outputs := make(map[string]interface{}, 0)
	deployment := response["deployment"].(map[string]interface{})
	resources := deployment["resources"].([]interface{})
	expected_urn := fmt.Sprintf("urn:pulumi:%s::%s::pulumi:pulumi:Stack::%s-%s", stack, project, project, stack)
	for _, resourceInterface := range resources {
		resource := resourceInterface.(map[string]interface{})

		log.Println(fmt.Sprintf("[DEBUG] FOUND URN: %s", resource["urn"].(string)))

		if resource["urn"].(string) == expected_urn {
			stack_outputs = resource["outputs"].(map[string]interface{})
		}
	}

	// Set outputs on the resource
	if err := d.Set("version", response["version"].(float64)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("stack_outputs", stack_outputs); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func dataSourceStackOutputs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStackOutputsRead,
		Schema: map[string]*schema.Schema{
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "organization name",
			},
			"project": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "project name",
			},
			"stack": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "stack name",
			},
			"version": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Version of Pulumi's API used",
			},
			"stack_outputs": &schema.Schema{
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Outputs for the Pulumi Stack",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
