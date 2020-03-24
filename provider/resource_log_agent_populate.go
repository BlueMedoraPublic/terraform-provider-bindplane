package provider

import (
	"strings"
	"time"

	"github.com/BlueMedoraPublic/bpcli/bindplane/sdk"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceLogAgentPopulate() *schema.Resource {
	return &schema.Resource{
		Create: resourceLogAgentCreate,
		Read:   resourceLogAgentRead,
		Update: resourceLogAgentUpdate,
		Delete: resourceLogAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			/*
				do not force new resource if name changes, this is
				okay as long as we already have the ID in the state
			*/
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			/*
				amount of time to wait for the agent to populate
				in the ui. resource creation fails if this is set
				too low
			*/
			"provisioning_timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
				Description: "time in seconds (30 or higher) to wait for a log agent to deploy",
			},

			/*
				values that could change on the api side,
				do not force replace
			*/
			"version": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: false,
			},
			"latest_version": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: false,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: false,
			},
		},
	}
}

func resourceLogAgentCreate(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	timeout := int32(d.Get("provisioning_timeout").(int))

	if timeout < 30 {
		return errors.New("timeout less than 30 seconds is not allowed for log agent creation")
	}

	a, err := waitForAgent(name, timeout)
	if err != nil {
		return err
	}

	d.SetId(a.ID)
	d.Set("version", a.Version)
	d.Set("latest_version", a.LatestVersion)
	d.Set("status", a.Status)
	return resourceLogAgentRead(d, m)
}

func resourceLogAgentRead(d *schema.ResourceData, m interface{}) error {
	a, err := bp.GetLogAgent(d.Id())
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "no agent with id") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", a.Name)
	d.Set("version", a.Version)
	d.Set("latest_version", a.LatestVersion)
	d.Set("status", a.Status)
	return nil
}

func resourceLogAgentUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceLogAgentRead(d, m)
}

func resourceLogAgentDelete(d *schema.ResourceData, m interface{}) error {
	if err := bp.DeleteLogAgent(d.Id()); err != nil {
		return err
	}
	return resourceLogAgentRead(d, m)
}

func waitForAgent(name string, timeout int32) (sdk.LogAgent, error) {
	start := int32(time.Now().Unix())
	for {
		current := int32(time.Now().Unix())
		if (current - start) > timeout {
			err := errors.New("timed out while waiting for agent to create: " + name)
			return sdk.LogAgent{}, err
		}

		a, err := bp.ListLogAgents()
		if err != nil {
			return sdk.LogAgent{}, err
		}

		if err := uniqueAgent(name, a); err != nil {
			return sdk.LogAgent{}, err
		}

		for _, a := range a {
			if a.Name == name {
				return a, nil
			}
		}

		// sleep for two seconds between API calls
		time.Sleep(2 * time.Second)
	}
}

func uniqueAgent(name string, a []sdk.LogAgent) error {
	c := 0
	for i := range a {
		if a[i].Name == name {
			c = c + 1
		}

		if c > 1 {
			return errors.New("agent " + name + " exists more than once in the API")
		}
	}
	return nil
}
