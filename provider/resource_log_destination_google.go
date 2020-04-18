package provider

import (
	"strings"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"

	"github.com/hashicorp/terraform/helper/schema"
)

const gcpDestTypeID = "stackdriver"

func resourceLogDestinationGCP() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogDestinationGCPCreate,
		Read:   resourceLogDestinationGCPRead,
		Delete: resourceLogDestinationGCPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"credentials": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"location": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "us-west1",
			},
			/*
			 destination_version is placed into the statefile however the provider
			 will not attempt to detect configuration drift on the API side.
			 if a version upgrade is desired, change it in your terraform config
			 and re-apply.
			*/
			"destination_version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLogDestinationGCPCreate(d *schema.ResourceData, m interface{}) error {
	config := sdk.LogDestConfig{
		Name:               d.Get("name").(string),
		DestinationTypeID:  gcpDestTypeID,
		DestinationVersion: d.Get("destination_version").(string),
		Configuration:      make(map[string]interface{}),
	}
	config.Configuration["credentials"] = d.Get("credentials").(string)
	config.Configuration["location"] = d.Get("location").(string)

	if config.DestinationVersion == "" {
		v, err := defaultLogDestVersion(gcpDestTypeID)
		if err != nil {
			return err
		}
		config.DestinationVersion = v
	}

	x, err := bp.CreateLogDestConfig(config)
	if err != nil {
		return err
	}

	d.SetId(x.ID)
	return resourceLogDestinationGCPRead(d, m)
}

func resourceLogDestinationGCPRead(d *schema.ResourceData, m interface{}) error {
	dest, err := bp.GetLogDestConfig(d.Id())
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

func resourceLogDestinationGCPDelete(d *schema.ResourceData, m interface{}) error {
	if err := bp.DelLogDestConfig(d.Id()); err != nil {
		return err
	}
	return resourceLogDestinationGCPRead(d, m)
}
