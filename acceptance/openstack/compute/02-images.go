package main

import (
	"fmt"
	"os"
	"github.com/rackspace/gophercloud/openstack/compute/images"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
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
	fmt.Fprintln(w, "ID\tRegion\tName\tStatus\tCreated\t")

	region := os.Getenv("OS_REGION_NAME")
	n := 0
	for _, ep := range eps {
		client := images.NewClient(ep.PublicURL, a, ao)

		if (region != "") && (region != ep.Region) {
			continue
		}

		listResults, err := images.List(client)
		if err != nil {
			panic(err)
		}

		imgs, err := images.GetImages(listResults)
		if err != nil {
			panic(err)
		}

		n = n + len(imgs)

		for _, i := range imgs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", i.Id, ep.Region, i.Name, i.Status, i.Created)
		}
	}
	w.Flush()
	fmt.Printf("--------\n%d images listed.\n", n)
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

