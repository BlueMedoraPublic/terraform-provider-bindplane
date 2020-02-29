package provider

import (
	"strings"

	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/logs/source"

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
			"configuration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceLogSourceCreate(d *schema.ResourceData, m interface{}) error {
	configuration := d.Get("configuration").(string)
	payload := []byte(configuration)
	x, err := source.Create(payload)
	if err != nil {
		return err
	}
	d.SetId(x)
	return resourceLogSourceRead(d, m)
}

func resourceLogSourceRead(d *schema.ResourceData, m interface{}) error {
	if err := source.Read(d.Id()); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no source config with id") {
			d.SetId("")
			return nil
		}
		return err
	}
	return nil
}

func resourceLogSourceDelete(d *schema.ResourceData, m interface{}) error {
	if err := source.Delete(d.Id()); err != nil {
		return err
	}

	/*
	 we do not remove the id here, instead we return the result
	 of resourceLogSourceRead() which will always set the id to
	 zero when it discovers that the resource no longer exists.
	 If we did it twice, resourceLogSourceRead() would perform a GET
	 against https:public-api.bindplane.bluemedora.com/v1/logs/source_configs/
     (because the id is empty!)	which returns a valid HTTP request (list all source configs)
	 which results in json unmarshal errors. This is not the APIs fault,
	 but a design choice for this source base. Hashicorp recomends
	 always returning the read function when performing create and
	 delete operations.
	*/
	//d.SetId("")
	return resourceLogSourceRead(d, m)
}
