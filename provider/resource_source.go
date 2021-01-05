package provider

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"reflect"
	"encoding/json"

	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/source"
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/util/jsonutil"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
	"github.com/hashicorp/terraform/helper/schema"
)

// sourceTimeoutMin is the minimum allowed timeout for a source
const sourceTimeoutMin = int(90)

func resourceSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceSourceCreate,
		Read:   resourceSourceRead,
		Delete: resourceSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"configuration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"collection_interval": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"collector_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"credential_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioning_timeout": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSourceCreate(d *schema.ResourceData, m interface{}) error {
	s, err := initSource(d)
	if err != nil {
		return err
	}

	timeout := timeout(d)

	if timeout < sourceTimeoutMin {
		return badTimeoutErr()
	}

	x, err := source.Create(bp, s, timeout)
	if err != nil {
		return err
	}

	d.SetId(x.SourceID)
	d.Set("job_id", x.JobID)
	return resourceSourceRead(d, m)
}

func resourceSourceRead(d *schema.ResourceData, m interface{}) error {
	const notFound = "Target could not be found"

	source, err := bp.GetSource(d.Id())
	if err != nil {
		if strings.Contains(err.Error(), notFound) {
			d.SetId("")
			return nil
		}
		return err
	}

	identical, err := sourceCompare(source, d.Get("configuration").(string))
	if err != nil {
		return err
	}

	if ! identical {
		d.Set("configuration", "")
	}

	d.Set("collection_interval", source.CollectionInterval)
	d.Set("collector_id", source.Collector.ID)
	d.Set("name", source.Name)
	d.Set("source_type", source.SourceType.ID)

	// TODO: Credentials should probably be an array, however, it is unlikely
	// that a source will have multiple credentials attached. This will require
	// a schema change, meaning the next major release of this provider will include
	// this improvement
	if len(source.Credentials) > 0 {
		d.Set("credential_id", source.Credentials[0].ID)
	} else {
		d.Set("credential_id", "")
	}

	return nil
}

func resourceSourceDelete(d *schema.ResourceData, m interface{}) error {
	if _, err := bp.DeleteSource(d.Id()); err != nil {
		return err
	}
	return resourceSourceRead(d, m)
}

func initSource(d *schema.ResourceData) (sdk.SourceConfigCreate, error) {
	s := sdk.SourceConfigCreate{}
	s.CollectionInterval = d.Get("collection_interval").(int)
	s.CollectorID = d.Get("collector_id").(string)
	s.Name = d.Get("name").(string)
	s.SourceType = d.Get("source_type").(string)
	s.Credentials.Credentials = d.Get("credential_id").(string)

	conf := d.Get("configuration").(string)
	if err := jsonutil.JSONToInterface(conf, &s.Configuration); err != nil {
		return s, err
	}

	return s, nil
}

func sourceCompare(remote sdk.SourceConfigGet, local string) (bool, error) {
	configuration := make(map[string]interface{})
	if err := json.Unmarshal([]byte(local), &configuration); err != nil {
		return false, err
	}
	return reflect.DeepEqual(remote.Configuration, configuration), nil
}

func timeout(d *schema.ResourceData) int {
	return d.Get("provisioning_timeout").(int)
}

func badTimeoutErr() error {
	return errors.New("BindPlane source timeout cannot be less than " + strconv.Itoa(sourceTimeoutMin) + " seconds")
}
