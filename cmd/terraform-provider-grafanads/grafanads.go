package main

import (
	"github.com/hashicorp/terraform/plugin"

	gprovider "github.com/jwierzbo/terraform-provider-grafana-datasource/pkg/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gprovider.Provider,
	})
}
