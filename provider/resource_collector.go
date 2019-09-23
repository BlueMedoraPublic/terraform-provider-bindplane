package provider

import (
	"github.com/pkg/errors"
	"strings"

	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/collector"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCollector() *schema.Resource {
	return &schema.Resource{
		Create: resourceCollectorCreate,
		Read:   resourceCollectorRead,
		Delete: resourceCollectorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

/*
resourceCollectorCreate scans the API for the collector and
maps it's ID to the resourceID. This resource does not deploy
a virtual machine, you should use your cloud provider's provider
with a userdata script for collector deployment
*/
func resourceCollectorCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	id, err := collector.WaitForAPI(name)
	if err != nil {
		return errors.Wrap(err, "terraform bindplane_collector resource failed while waiting for collector to appear in bindplane api")
	}

	d.Set("group", id)
	d.SetId(id)
	return nil
}

// resourceCollectorRead checks to see if a specific collector exists
func resourceCollectorRead(d *schema.ResourceData, m interface{}) error {
	if err := collector.Read(d.Id()); err != nil {
		/*
			It is possible the collector in the tf state does not exist.
			If this happens, remove it from the tf state by setting
			it's id to ""
		*/
		if strings.Contains(err.Error(), "could not be found") {
			d.SetId("")
			return nil
		}
		return err
	}
	return nil
}

/*
resourceCollectorDelete removes the collector from the BindPlane
API. This resource does not delete your virtual machine, just like
the create resource does not create a virtual machine
*/
func resourceCollectorDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Get("group").(string) + "/" + d.Id()
	if err := collector.Delete(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
