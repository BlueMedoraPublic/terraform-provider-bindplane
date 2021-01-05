package provider

import (
	"fmt"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

var bp *sdk.BindPlane

const envAPIKey = "BINDPLANE_API_KEY"

// Provider is the Bindplane Terraform Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:     schema.TypeString,
				Optional: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"BINDPLANE_API_KEY",
				}, nil),
				ValidateFunc: validUUID,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bindplane_credential":             resourceCredential(),
			"bindplane_source":                 resourceSource(),
			"bindplane_collector":              resourceCollector(),
			"bindplane_log_source":             resourceLogSource(),
			"bindplane_log_destination":        resourceLogDestination(),
			"bindplane_log_destination_google": resourceLogDestinationGCP(),
			"bindplane_log_template":           resourceLogTemplate(),
			"bindplane_log_agent_populate":     resourceLogAgentPopulate(),
			"bindplane_log_bind_source":        resourceLogBindSource(),
			"bindplane_log_bind_destination":   resourceLogBindDestination(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"bindplane_agent_install_cmd": dataSourceAgentInstallCMD(),
		},
		ConfigureFunc: initBindplane,
	}
}

func initBindplane(d *schema.ResourceData) (interface{}, error) {
	bp = new(sdk.BindPlane)

	if d.Get("api_key").(string) == "" {
		return d, apiKeyRequiredErr()
	}

	if err := bp.Init(); err != nil {
		return d, err
	}

	return d, testConnection(bp)
}

func testConnection(bp *sdk.BindPlane) error {
	if _, err := bp.ListJobs(); err != nil {
		return errors.Wrap(err, "Test connection failed, could not list jobs. Is the API key correct?")
	}
	return nil
}

func apiKeyRequiredErr() error {
	return errors.New("BindPlane API Key is not set. You must set the API key in " +
		"one of the following: BindPlane provider config, BINDPLANE_API_KEY environment " +
		"variable, or configure '.bpcli' in your home directory.")
}

func validUUID(v interface{}, k string) (warnings []string, errors []error) {
	if v == nil || v.(string) == "" {
		return
	}
	value := v.(string)
	if _, err := uuid.Parse(value); err != nil {
		errors = append(errors,
			fmt.Errorf(k+" is not a valid uuid: ", value, err))
	}
	return
}
