package provider

import (
	"os"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/BlueMedoraPublic/bpcli/util/uuid"
	"github.com/pkg/errors"
)

var bp *sdk.BindPlane

// Provider is the Bindplane Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"bindplane_credential":           resourceCredential(),
			"bindplane_source":               resourceSource(),
			"bindplane_collector":            resourceCollector(),
			"bindplane_log_source":           resourceLogSource(),
			"bindplane_log_destination":      resourceLogDestination(),
			"bindplane_log_template":         resourceLogTemplate(),
			"bindplane_log_agent_populate":   resourceLogAgentPopulate(),
			"bindplane_log_bind_source":      resourceLogBindSource(),
			"bindplane_log_bind_destination": resourceLogBindDestination(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bindplane_agent_install_cmd": dataSourceAgentInstallCMD(),
		},
		ConfigureFunc: initBindplane,
	}
}

func initBindplane(d *schema.ResourceData) (interface{}, error) {
	if err := checkEnv(); err != nil {
		return nil, errors.Wrap(err, "Not attempting to initilize the bindplane sdk")
	}

	bp = new(sdk.BindPlane)
	if err := bp.Init(); err != nil {
		return d, err
	}

	return d, testConnection(bp)
}

func checkEnv() error {
	apiKey := os.Getenv("BINDPLANE_API_KEY")
	if len(apiKey) == 0 {
		return errors.New("required environment variable BINDPLANE_API_KEY is not set")
	}

	if uuid.IsUUID(apiKey) == false {
		return errors.New("required environment variable BINDPLANE_API_KEY is set but does not appear to be a uuid")
	}

	return nil
}

func testConnection(bp *sdk.BindPlane) error {
	if _, err := bp.ListJobs(); err != nil {
		return errors.Wrap(err, "Test connection failed, could not list jobs. Is the API key correct?")
	}
	return nil
}
