package main

import (
	"fmt"
	"github.com/rackspace/gophercloud/openstack/compute/servers"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
)

func main() {
	ao, err := utils.AuthOptions()
	if err != nil {
		panic(err)
	}

	a, err := identity.Authenticate(ao)
	if err != nil {
		panic(err)
	}

	sc, err := identity.GetServiceCatalog(a)
	if err != nil {
		panic(err)
	}

	ep, err := findAnyComputeEndpoint(sc)
	if err != nil {
		panic(err)
	}

	client := servers.NewClient(ep, a, ao)

	listResults, err := servers.List(client)
	if err != nil {
		panic(err)
	}

	svrs, err := servers.GetServers(listResults)
	if err != nil {
		panic(err)
	}

	for _, s := range svrs {
		fmt.Printf("ID(%s)\n", s.Id)
		fmt.Printf("    Name(%s)\n", s.Name)
		fmt.Printf("    IPv4(%s)\n    IPv6(%s)\n", s.AccessIPv4, s.AccessIPv6)
	}
	fmt.Printf("--------\n%d servers listed.\n", len(svrs))
}


func findAnyComputeEndpoint(sc *identity.ServiceCatalog) (string, error) {
	ces, err := sc.CatalogEntries()
	if err != nil {
		return "", err
	}

	for _, ce := range ces {
		if ce.Type == "compute" {
			return ce.Endpoints[0].PublicURL, nil
		}
	}

	return "", fmt.Errorf("Compute endpoint not found.")
}

