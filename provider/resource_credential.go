package provider

import (
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCredential() *schema.Resource {
	return &schema.Resource{
		Create: resourceCredentialCreate,
		Read:   resourceCredentialRead,
		Delete: resourceCredentialDelete,
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

func resourceCredentialCreate(d *schema.ResourceData, m interface{}) error {
	configuration := d.Get("configuration").(string)
	payload := []byte(configuration)
	x, err := bp.CreateCredential(payload)
	if err != nil {
		return err
	}
	d.SetId(x.ID)
	return resourceCredentialRead(d, m)
}

func resourceCredentialRead(d *schema.ResourceData, m interface{}) error {
	if _, err := bp.GetCredential(d.Id()); err != nil {
		/*
			It is possible the credential in the tf state does not exist.
			If this happens, remove it from the tf state by setting
			it's id to ""
		*/
		if strings.Contains(err.Error(), "No credential with ID") {
			d.SetId("")
			return nil
		}
		return err
	}

	// dont set anything
	return nil
}

func resourceCredentialDelete(d *schema.ResourceData, m interface{}) error {
	if err := bp.DeleteCredential(d.Id()); err != nil {
		return err
	}
	return resourceCredentialRead(d, m)
}
