package provider

import (
	"github.com/pkg/errors"
	"strconv"
	"strings"

	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/source"
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/util/jsonutil"
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/util/trim"

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

	x, err := source.Create(s, timeout)
	if err != nil {
		return err
	}

	d.SetId(x.SourceID)
	d.Set("job_id", x.JobID)
	return resourceSourceRead(d, m)
}

func resourceSourceRead(d *schema.ResourceData, m interface{}) error {
	source, err := source.Read(d.Id())
	if err != nil {
		// remove from state if not exist
		if strings.Contains(err.Error(), "Target could not be found") {
			d.SetId("")
			return nil
		}
		return err
	}

	/*
	  if the state source configuration is different from
	  the configuration returned by the API, force the source
	  to be rebuild
	*/
	c, err := confDiff(d, source)
	if err != nil {
		return err
	} else if c == false {
		d.Set("configuration", "")
	}

	d.Set("collection_interval", source.CollectionInterval)
	d.Set("collector_id", source.Collector.ID)
	d.Set("name", source.Name)
	d.Set("source_type", source.SourceType.ID)

	// some sources will not havea a credential, which will
	// cause terraform to panic if we skip this check
	if len(source.Credentials) > 0 {
		d.Set("credential_id", source.Credentials[0].ID)
	} else {
		d.Set("credential_id", "")
	}

	return nil
}

func resourceSourceDelete(d *schema.ResourceData, m interface{}) error {
	if err := source.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
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

// returns true if api and state source configurations are equal
func confDiff(d *schema.ResourceData, s sdk.SourceConfigGet) (bool, error) {
	// api response as bytes
	apiConf, err := jsonutil.InterfaceToJSONBytes(s.Configuration)
	if err != nil {
		return false, err
	}

	// force the source to rebuild if state differs from api
	// response
	stateConf, err := trim.Trim(d.Get("configuration").(string))
	if string(apiConf) != stateConf {
		return false, nil
	}
	return true, nil
}

func timeout(d *schema.ResourceData) int {
	return d.Get("provisioning_timeout").(int)
}

func badTimeoutErr() error {
	return errors.New("BindPlane source timeout cannot be less than " + strconv.Itoa(sourceTimeoutMin) + " seconds")
}
