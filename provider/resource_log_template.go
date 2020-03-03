package provider

import (
	"strings"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/logs/template"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLogTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogTemplateCreate,
		Read:   resourceLogTemplateRead,
		Delete: resourceLogTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_config_ids": {
				// source_config_ids is a list and looks like this:
				// source_type_ids = [a, b, c]
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				ForceNew: true,
			},
			"destination_config_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"agent_group": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLogTemplateCreate(d *schema.ResourceData, m interface{}) error {
	t := sdk.LogTemplate{}
	t.Name = d.Get("name").(string)
	t.DestinationConfigID = d.Get("destination_config_id").(string)
	t.AgentGroup = d.Get("agent_group").(string)

	if v := d.Get("source_config_ids"); v != nil {
		vs := v.(*schema.Set)
		t.SourceConfigIds = make([]string, vs.Len())
		for i, v := range vs.List() {
			t.SourceConfigIds[i] = v.(string)
		}
	}

	x, err := template.Create(t)
	if err != nil {
		return err
	}

	d.SetId(x)
	return resourceLogTemplateRead(d, m)
}

func resourceLogTemplateRead(d *schema.ResourceData, m interface{}) error {
	t, err := template.Read(d.Id())
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no template with id") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", t.Name)
	d.Set("destination_config_id", t.DestinationConfigID)
	d.Set("agent_group", t.AgentGroup)

	a := make([]interface{}, len(t.SourceConfigIds))
	for i, str := range t.SourceConfigIds {
		a[i] = str
	}
	d.Set("source_config_ids", a)

	return nil
}

func resourceLogTemplateDelete(d *schema.ResourceData, m interface{}) error {
	if err := template.Delete(d.Id()); err != nil {
		return err
	}
	return resourceLogTemplateRead(d, m)
}
