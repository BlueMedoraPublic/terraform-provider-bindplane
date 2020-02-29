package provider

import (
	"strings"

	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/logs/destination"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLogDestination() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogDestinationCreate,
		Read:   resourceLogDestinationRead,
		Delete: resourceLogDestinationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLogDestinationCreate(d *schema.ResourceData, m interface{}) error {
	configuration := d.Get("configuration").(string)
	payload := []byte(configuration)
	x, err := destination.Create(payload)
	if err != nil {
		return err
	}
	d.SetId(x)
	return resourceLogDestinationRead(d, m)
}

func resourceLogDestinationRead(d *schema.ResourceData, m interface{}) error {
	if err := destination.Read(d.Id()); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no destination config with id") {
			d.SetId("")
			return nil
		}
		return err
	}
	return nil
}

func resourceLogDestinationDelete(d *schema.ResourceData, m interface{}) error {
	if err := destination.Delete(d.Id()); err != nil {
		return err
	}
	return resourceLogDestinationRead(d, m)
}
