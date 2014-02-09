package main

import (
	"fmt"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
)

type extractor func(*identity.Token) string

func main() {
	// Create an initialized set of authentication options based on available OS_*
	// environment variables.
	ao, err := utils.AuthOptions()
	if err != nil {
		panic(err)
	}

	// Attempt to authenticate with them.
	r, err := identity.Authenticate(ao)
	if err != nil {
		panic(err)
	}

	// We're authenticated; now let's grab our authentication token.
	t, err := identity.GetToken(r)
	if err != nil {
		panic(err)
	}

	// Authentication tokens have a variety of fields which might be of some interest.
	// Let's print a few of them out.
	table := map[string]extractor{
		"ID":      func(t *identity.Token) string { return t.Id },
		"Expires": func(t *identity.Token) string { return t.Expires },
	}

	for attr, fn := range table {
		fmt.Printf("Your token's %s is %s\n", attr, fn(t))
	}

	// With each authentication, you receive a master directory of all the services
	// your account can access.  This "service catalog", as OpenStack calls it,
	// provides you the means to exploit other OpenStack services.
	sc, err := identity.GetServiceCatalog(r)
	if err != nil {
		panic(err)
	}

	// Different providers will provide different services.  Let's print them
	// in summary.
	ces, err := sc.CatalogEntries()
	fmt.Printf("Service Catalog Summary:\n  %32s   %-16s\n", "Name", "Type")
	for _, ce := range ces {
		fmt.Printf("  %32s | %-16s\n", ce.Name, ce.Type)
	}

	// Now let's print them in greater detail.
	for _, ce := range ces {
		fmt.Printf("Endpoints for %s/%s\n", ce.Name, ce.Type)
		for _, ep := range ce.Endpoints {
			fmt.Printf("  Version: %s\n", ep.VersionId)
			fmt.Printf("  Region: %s\n", ep.Region)
			fmt.Printf("  Tenant: %s\n", ep.TenantId)
			fmt.Printf("  Public URL: %s\n", ep.PublicURL)
			fmt.Printf("  Internal URL: %s\n", ep.InternalURL)
		}
	}
}
