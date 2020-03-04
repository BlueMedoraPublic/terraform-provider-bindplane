package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func dataSourceAgentInstallCMD() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAgentInstallCMDRead,
		Schema: map[string]*schema.Schema{
			"platform": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"command": {
				Type:     schema.TypeString,
				Computed: true,
				// we do not want to force a rebuild if the
				// command changes!
				ForceNew: false,
			},
		},
	}
}

func dataSourceAgentInstallCMDRead(d *schema.ResourceData, m interface{}) error {
	platform := d.Get("platform").(string)

	c, err := installCmd(platform)
	if err != nil {
		return errors.Wrap(err, "cannot retrieve data source 'bindplane_agent_install_cmd'")
	}

	d.Set("command", c)
	d.SetId(platform + "=" + c)
	return nil
}

func installCmd(platform string) (string, error) {
	if platform == "all" {
		return "", errors.New("platform 'all' cannot be used, specify a specific platform")
	}
	return bp.InstallCMDLogAgent(platform)
}
