package provider

import (
	"encoding/json"
	"strings"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/logs/source"
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/util/compare"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLogSource() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogSourceCreate,
		Read:   resourceLogSourceRead,
		Delete: resourceLogSourceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_type_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			/*
			 source version is placed into the statefile however the provider
			 will not attempt to detect configuration drift on the API side.
			 if a version upgrade is desired, change it in your terraform code
			 and re-apply.
			*/
			"source_version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			/*
				 leaving custom_template un implemented as I believe it is
				 not needed with terraform. If enabled, make sure to to add
				 d.Set("source_version", s.Source.Version) to resourceLogSourceRead
				"custom_template": {
					Type: schema.TypeString,
					Optional: true,
					ForceNew: true,
				},*/
			"configuration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLogSourceCreate(d *schema.ResourceData, m interface{}) error {
	config := sdk.LogSourceConfig{
		Name:          d.Get("name").(string),
		SourceTypeID:  d.Get("source_type_id").(string),
		SourceVersion: d.Get("source_version").(string),
		Configuration: make(map[string]interface{}),
	}

	c := []byte(d.Get("configuration").(string))
	if err := json.Unmarshal(c, &config.Configuration); err != nil {
		return err
	}

	x, err := source.Create(config)
	if err != nil {
		return err
	}

	d.SetId(x)
	return resourceLogSourceRead(d, m)
}

func resourceLogSourceRead(d *schema.ResourceData, m interface{}) error {
	s, err := source.Read(d.Id())
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no source config with id") {
			d.SetId("")
			return nil
		}
		return err
	}

	same, err := logSourceConfigDiff(d, s)
	if err != nil {
		return err
	}

	// if state differs from api, unset local copy to force resource
	// replacement
	if !same {
		d.Set("configuration", "")
	}

	d.Set("name", s.Name)
	d.Set("source_type_id", s.Source.ID)

	return nil
}

func resourceLogSourceDelete(d *schema.ResourceData, m interface{}) error {
	if err := source.Delete(d.Id()); err != nil {
		return err
	}
	return resourceLogSourceRead(d, m)
}

func logSourceConfigDiff(d *schema.ResourceData, apiConf sdk.LogSourceConfig) (bool, error) {
	stateConfBytes := []byte(d.Get("configuration").(string))
	if len(stateConfBytes) == 0 {
		return false, nil
	}

	stateConf := sdk.LogSourceConfig{}
	if err := json.Unmarshal(stateConfBytes, &stateConf.Configuration); err != nil {
		return false, err
	}

	return compare.MapStringInterface(stateConf.Configuration, apiConf.Configuration), nil
}
