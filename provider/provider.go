package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider is the Bindplane Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"bindplane_credential": resourceCredential(),
			"bindplane_source":     resourceSource(),
			"bindplane_collector":  resourceCollector(),
			"bindplane_log_source": resourceLogSource(),
		},
	}
}
