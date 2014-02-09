package main

import (
	"fmt"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
)

func main() {
	// Create an initialized set of authentication options based on available OS_*
	// environment variables.
	ao, err := utils.AuthOptions()
	if err != nil {
		panic(err)
	}

	// Attempt to query extensions.
	exts, err := identity.GetExtensions(ao)
	if err != nil {
		panic(err)
	}

	// Print out a summary of supported extensions
	aliases, err := exts.Aliases()
	if err != nil {
		panic(err)
	}
	fmt.Println("Extension Aliases:")
	for _, alias := range aliases {
		fmt.Printf("  %s\n", alias)
	}
}
