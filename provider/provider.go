package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		DataSourcesMap: map[string]*schema.Resource{
			"hello_world": dataSourceWorld(),
		},
	}
}

func dataSourceWorld() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWorldRead,
		Schema: map[string]*schema.Schema{
			"message": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceWorldRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId("world")
	d.Set("message", "Hello, World!")
	return nil
}
