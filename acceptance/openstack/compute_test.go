package openstack

import (
	"testing"
	"fmt"
	"github.com/rackspace/gophercloud/openstack/compute/servers"
	"github.com/rackspace/gophercloud/openstack/compute/images"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
	"os"
	"text/tabwriter"
)

func TestListServers(t *testing.T) {
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Error(err)
		return
	}

	a, err := identity.Authenticate(ao)
	if err != nil {
		t.Error(err)
		return
	}

	sc, err := identity.GetServiceCatalog(a)
	if err != nil {
		t.Error(err)
		return
	}

	eps, err := findAllComputeEndpoints(sc)
	if err != nil {
		t.Error(err)
		return
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
			t.Error(err)
			return
		}

		svrs, err := servers.GetServers(listResults)
		if err != nil {
			t.Error(err)
			return
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

func TestListImages(t *testing.T) {
	ao, err := utils.AuthOptions()
	if err != nil {
		t.Error(err)
		return
	}

	a, err := identity.Authenticate(ao)
	if err != nil {
		t.Error(err)
		return
	}

	sc, err := identity.GetServiceCatalog(a)
	if err != nil {
		t.Error(err)
		return
	}

	eps, err := findAllComputeEndpoints(sc)
	if err != nil {
		t.Error(err)
		return
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 2, 8, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tRegion\tName\tStatus\tCreated\t")

	region := os.Getenv("OS_REGION_NAME")
	n := 0
	for _, ep := range eps {
		if (region != "") && (region != ep.Region) {
			continue
		}

		client := images.NewClient(ep.PublicURL, a, ao)

		listResults, err := images.List(client)
		if err != nil {
			t.Error(err)
			return
		}

		imgs, err := images.GetImages(listResults)
		if err != nil {
			t.Error(err)
			return
		}

		n = n + len(imgs)

		for _, i := range imgs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", i.Id, ep.Region, i.Name, i.Status, i.Created)
		}
	}
	w.Flush()
	fmt.Printf("--------\n%d images listed.\n", n)
}

