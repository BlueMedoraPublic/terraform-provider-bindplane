package provider

import (
	"strings"
	"encoding/json"

	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/logs/destination"
	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"

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
			"name": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"destination_type_id": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			/*
			 destination_version is placed into the statefile however the provider
			 will not attempt to detect configuration drift on the API side.
			 if a version upgrade is desired, change it in your terraform config
			 and re-apply.
			*/
			"destination_version": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"configuration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLogDestinationCreate(d *schema.ResourceData, m interface{}) error {
	config := sdk.LogDestConfig{
		Name:               d.Get("name").(string),
		DestinationTypeID:  d.Get("destination_type_id").(string),
		DestinationVersion: d.Get("destination_version").(string),
		Configuration: make(map[string]interface{}),
	}

	c := []byte(d.Get("configuration").(string))
	if err := json.Unmarshal(c, &config.Configuration); err != nil {
		return err
	}

	x, err := destination.Create(config)
	if err != nil {
		return err
	}

	d.SetId(x)
	return resourceLogDestinationRead(d, m)
}

func resourceLogDestinationRead(d *schema.ResourceData, m interface{}) error {
	dest, err := destination.Read(d.Id())
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no destination config with id") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", dest.Name)
	d.Set("destination_type_id", dest.Destination.ID)

	return nil
}

func resourceLogDestinationDelete(d *schema.ResourceData, m interface{}) error {
	if err := destination.Delete(d.Id()); err != nil {
		return err
	}
	return resourceLogDestinationRead(d, m)
}
