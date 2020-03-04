package provider

import (
	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/common"

	"github.com/hashicorp/terraform/helper/schema"
)

var bp *sdk.BindPlane

// Provider is the Bindplane Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"bindplane_credential":      resourceCredential(),
			"bindplane_source":          resourceSource(),
			"bindplane_collector":       resourceCollector(),
			"bindplane_log_source":      resourceLogSource(),
			"bindplane_log_destination": resourceLogDestination(),
			"bindplane_log_template":    resourceLogTemplate(),
			"bindplane_log_agent_populate":       resourceLogAgentPopulate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bindplane_agent_install_cmd": dataSourceAgentInstallCMD(),
		},
		ConfigureFunc: initBindplane,
	}
}

func initBindplane(d *schema.ResourceData) (interface{}, error) {
	var err error
	bp, err = common.New()
	return d, err
}
