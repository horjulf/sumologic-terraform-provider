package main

import (
	"github.com/SumoLogic/sumologic-terraform-provider/sumologic"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: sumologic.Provider,
	})
}
