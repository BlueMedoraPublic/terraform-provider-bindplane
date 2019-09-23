package provider

import (
	"strings"

	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider/bindplane/credential"

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
	x, err := credential.Create(payload)
	if err != nil {
		return err
	}
	d.SetId(x)
	return nil
}

func resourceCredentialRead(d *schema.ResourceData, m interface{}) error {
	if err := credential.Read(d.Id()); err != nil {
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
	if err := credential.Delete(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
