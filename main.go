package main

import (
	"github.com/BlueMedoraPublic/terraform-provider-bindplane/provider"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

const version = "0.2.0-beta.1"

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return provider.Provider()
		},
	})
}
