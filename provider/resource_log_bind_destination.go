package provider

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLogBindDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogBindDestCreate,
		Read:   resourceLogBindDestRead,
		Delete: resourceLogBindDestDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"destination_config_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"agent_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLogBindDestCreate(d *schema.ResourceData, m interface{}) error {
	destID := d.Get("destination_config_id").(string)
	agentID  := d.Get("agent_id").(string)

	x, err := bp.DeployLogAgentDest(agentID, destID)
	if err != nil {
		return err
	}

	d.SetId(x.DestinationConfigID)
	return resourceLogBindDestRead(d, m)
}

func resourceLogBindDestRead(d *schema.ResourceData, m interface{}) error {
	destID := d.Id()
	agentID  := d.Get("agent_id").(string)

	t, err := bp.GetLogAgentDest(agentID, destID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "was not found when") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("destination_config_id", t.DestinationConfigID)
	d.Set("agent_id", agentID)
	return nil
}

func resourceLogBindDestDelete(d *schema.ResourceData, m interface{}) error {
	destID := d.Id()
	agentID  := d.Get("agent_id").(string)

	if err := bp.DeleteLogAgentDest(agentID, destID); err != nil {
		return err
	}
	return resourceLogBindDestRead(d, m)
}
