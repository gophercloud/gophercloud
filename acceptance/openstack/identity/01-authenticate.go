package main

import (
	"fmt"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
)

type extractor func(*identity.TokenDesc) string

func main() {
	ao, err := utils.AuthOptions()
	if err != nil {
		panic(err)
	}

	ao.AllowReauth = true
	r, err := identity.Authenticate(ao)
	if err != nil {
		panic(err)
	}

	t, err := identity.Token(r)
	if err != nil {
		panic(err)
	}

	table := map[string]extractor{
		"ID":      func(t *identity.TokenDesc) string { return t.Id() },
		"Expires": func(t *identity.TokenDesc) string { return t.Expires() },
	}

	for attr, fn := range table {
		fmt.Printf("Your token's %s is %s\n", attr, fn(t))
	}

	sc, err := identity.ServiceCatalog(r)
	if err != nil {
		panic(err)
	}
	ces, err := sc.CatalogEntries()
	fmt.Printf("Service Catalog Summary:\n  %32s   %-16s\n", "Name", "Type")
	for _, ce := range ces {
		fmt.Printf("  %32s | %-16s\n", ce.Name, ce.Type)
	}

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
