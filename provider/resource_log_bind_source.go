package provider

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLogBindSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogBindSourceCreate,
		Read:   resourceLogBindSourceRead,
		Delete: resourceLogBindSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"source_config_id": {
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

func resourceLogBindSourceCreate(d *schema.ResourceData, m interface{}) error {
	sourceID := d.Get("source_config_id").(string)
	agentID  := d.Get("agent_id").(string)

	x, err := bp.DeployLogAgentSource(agentID, sourceID)
	if err != nil {
		return err
	}

	d.SetId(x.SourceConfigID)
	return resourceLogBindSourceRead(d, m)
}

func resourceLogBindSourceRead(d *schema.ResourceData, m interface{}) error {
	sourceID := d.Id()
	agentID  := d.Get("agent_id").(string)

	t, err := bp.GetLogAgentSource(agentID, sourceID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "was not found when") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("source_config_id", t.SourceConfigID)
	d.Set("agent_id", agentID)
	return nil
}

func resourceLogBindSourceDelete(d *schema.ResourceData, m interface{}) error {
	sourceID := d.Id()
	agentID  := d.Get("agent_id").(string)

	if err := bp.DeleteLogAgentSource(agentID, sourceID); err != nil {
		return err
	}
	return resourceLogBindSourceRead(d, m)
}
