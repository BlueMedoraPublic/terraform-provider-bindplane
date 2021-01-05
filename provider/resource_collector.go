package provider

import (
	"github.com/pkg/errors"
	"strings"
	"time"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

// Max time to wait for a collector to deploy in seconds
const collectorTimeout = int64(300)

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
	id, err := waitForAPI(bp, name)
	if err != nil {
		return errors.Wrap(err, "terraform bindplane_collector resource failed while waiting for collector to appear in bindplane api")
	}

	d.Set("group", id)
	d.SetId(id)
	return resourceCollectorRead(d, m)
}

// resourceCollectorRead checks to see if a specific collector exists
func resourceCollectorRead(d *schema.ResourceData, m interface{}) error {
	const notFound = "could not be found"

	collector, err := bp.GetCollector(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), notFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("group", collector.ID)
	d.Set("name", collector.Name)

	return nil
}

/*
resourceCollectorDelete removes the collector from the BindPlane
API. This resource does not delete your virtual machine, just like
the create resource does not create a virtual machine
*/
func resourceCollectorDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Get("group").(string) + "/" + d.Id()
	if err := bp.DeleteCollector(id); err != nil {
		return err
	}
	return resourceCollectorRead(d, m)
}

/*
WaitForAPI checks the BindPlane API every ten seconds
until a collector with the name 'collectorName' is found. The
collector id is returned.
*/
func waitForAPI(bp *sdk.BindPlane, collectorName string) (string, error) {
	startTime := time.Now().Unix()

	for {
		collectors, err := bp.GetCollectors()
		if err != nil {
			return "", errors.Wrap(err, "sdk.GetCollectors() returned an error")
		}

		for _, collector := range collectors {
			if collector.Name == collectorName {
				return collector.ID, nil
			}
		}

		time.Sleep(10 * time.Second)
		timeCurrent := time.Now().Unix()
		if (timeCurrent - startTime) > collectorTimeout {
			msg := "Timeout exceeded for collector creation: " + "CollectorName: " + collectorName
			return "", errors.New(msg)
		}
	}
}
