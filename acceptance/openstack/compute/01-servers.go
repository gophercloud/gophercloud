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

	eps, err := findAllComputeEndpoints(sc)
	if err != nil {
		panic(err)
	}

	clients := make([]*servers.Client, len(eps))
	for i, ep := range eps {
		clients[i] = servers.NewClient(ep, a, ao)
	}

	n := 0
	for _, client := range clients {
		listResults, err := servers.List(client)
		if err != nil {
			panic(err)
		}

		svrs, err := servers.GetServers(listResults)
		if err != nil {
			panic(err)
		}

		n = n + len(svrs)

		for _, s := range svrs {
			fmt.Printf("ID(%s)\n", s.Id)
			fmt.Printf("    Name(%s)\n", s.Name)
			fmt.Printf("    IPv4(%s)\n    IPv6(%s)\n", s.AccessIPv4, s.AccessIPv6)
		}
	}
	fmt.Printf("--------\n%d servers listed.\n", n)
}


func findAllComputeEndpoints(sc *identity.ServiceCatalog) ([]string, error) {
	var eps []string

	ces, err := sc.CatalogEntries()
	if err != nil {
		return eps, err
	}

	for _, ce := range ces {
		if ce.Type == "compute" {
			eps := make([]string, len(ce.Endpoints))
			for i, endpoint := range ce.Endpoints {
				eps[i] = endpoint.PublicURL
			}
			return eps, nil
		}
	}

	return eps, fmt.Errorf("Compute endpoint not found.")
}

