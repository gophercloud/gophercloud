package main

import (
	"fmt"
	"github.com/rackspace/gophercloud/openstack/compute/servers"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
	"os"
	"text/tabwriter"
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

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 2, 8, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tRegion\tName\tIPv4\tIPv6\t")

	region := os.Getenv("OS_REGION_NAME")
	n := 0
	for _, ep := range eps {
		if (region != "") && (region != ep.Region) {
			continue
		}

		client := servers.NewClient(ep.PublicURL, a, ao)

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
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", s.Id, s.Name, ep.Region, s.AccessIPv4, s.AccessIPv6)
		}
	}
	w.Flush()
	fmt.Printf("--------\n%d servers listed.\n", n)
}

func findAllComputeEndpoints(sc *identity.ServiceCatalog) ([]identity.Endpoint, error) {
	ces, err := sc.CatalogEntries()
	if err != nil {
		return nil, err
	}

	for _, ce := range ces {
		if ce.Type == "compute" {
			return ce.Endpoints, nil
		}
	}

	return nil, fmt.Errorf("Compute endpoint not found.")
}
