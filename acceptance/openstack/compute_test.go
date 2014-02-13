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

type testState struct {
	o	identity.AuthOptions
	a	identity.AuthResults
	sc	*identity.ServiceCatalog
	eps	[]identity.Endpoint
	w	*tabwriter.Writer
}

func prepForTest() (*testState, error) {
	var err error

	ts := new(testState)

	ts.o, err = utils.AuthOptions()
	if err != nil {
		return ts, err
	}

	ts.a, err = identity.Authenticate(ts.o)
	if err != nil {
		return ts, err
	}

	ts.sc, err = identity.GetServiceCatalog(ts.a)
	if err != nil {
		return ts, err
	}

	ts.eps, err = findAllComputeEndpoints(ts.sc)
	if err != nil {
		return ts, err
	}

	ts.w = new(tabwriter.Writer)
	ts.w.Init(os.Stdout, 2, 8, 2, ' ', 0)

	return ts, nil
}

func TestListServers(t *testing.T) {
	ts, err := prepForTest()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Fprintln(ts.w, "ID\tRegion\tName\tIPv4\tIPv6\t")

	region := os.Getenv("OS_REGION_NAME")
	n := 0
	for _, ep := range ts.eps {
		if (region != "") && (region != ep.Region) {
			continue
		}

		client := servers.NewClient(ep.PublicURL, ts.a, ts.o)

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
			fmt.Fprintf(ts.w, "%s\t%s\t%s\t%s\t%s\t\n", s.Id, s.Name, ep.Region, s.AccessIPv4, s.AccessIPv6)
		}
	}
	ts.w.Flush()
	fmt.Printf("--------\n%d servers listed.\n", n)
}

func TestListImages(t *testing.T) {
	ts, err := prepForTest()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Fprintln(ts.w, "ID\tRegion\tName\tStatus\tCreated\t")

	region := os.Getenv("OS_REGION_NAME")
	n := 0
	for _, ep := range ts.eps {
		if (region != "") && (region != ep.Region) {
			continue
		}

		client := images.NewClient(ep.PublicURL, ts.a, ts.o)

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
			fmt.Fprintf(ts.w, "%s\t%s\t%s\t%s\t%s\t\n", i.Id, ep.Region, i.Name, i.Status, i.Created)
		}
	}
	ts.w.Flush()
	fmt.Printf("--------\n%d images listed.\n", n)
}

func TestListFlavors(t *testing.T) {
	ts, err := prepForTest()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Fprintln(ts.w, "ID\tRegion\tName\tStatus\tCreated\t")

	region := os.Getenv("OS_REGION_NAME")
	n := 0
	for _, ep := range ts.eps {
		if (region != "") && (region != ep.Region) {
			continue
		}

		client := flavors.NewClient(ep.PublicURL, ts.a, ts.o)

		listResults, err := flavors.List(client)
		if err != nil {
			t.Error(err)
			return
		}

		flavs, err := flavors.GetFlavors(listResults)
		if err != nil {
			t.Error(err)
			return
		}

		n = n + len(flavs)

		for _, f := range flavs {
			fmt.Fprintf(ts.w, "%s\t%s\t%s\t%s\t%s\t\n", f.Id, ep.Region, f.Name, f.Status, f.Created)
		}
	}
	ts.w.Flush()
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

