package openstack

import (
	"testing"
	"fmt"
	"github.com/rackspace/gophercloud/openstack/compute/servers"
	"github.com/rackspace/gophercloud/openstack/compute/images"
	"github.com/rackspace/gophercloud/openstack/compute/flavors"
	"github.com/rackspace/gophercloud/openstack/identity"
	"github.com/rackspace/gophercloud/openstack/utils"
	"os"
	"text/tabwriter"
	"time"
	"crypto/rand"
)

type testState struct {
	o	identity.AuthOptions
	a	identity.AuthResults
	sc	*identity.ServiceCatalog
	eps	[]identity.Endpoint
	w	*tabwriter.Writer
}

func setupForList() (*testState, error) {
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
	ts, err := setupForList()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Fprintln(ts.w, "ID\tRegion\tName\tStatus\tIPv4\tIPv6\t")

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
			fmt.Fprintf(ts.w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", s.Id, s.Name, ep.Region, s.Status, s.AccessIPv4, s.AccessIPv6)
		}
	}
	ts.w.Flush()
	fmt.Printf("--------\n%d servers listed.\n", n)
}

func TestListImages(t *testing.T) {
	ts, err := setupForList()
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
	ts, err := setupForList()
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Fprintln(ts.w, "ID\tRegion\tName\tRAM\tDisk\tVCPUs\t")

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
			fmt.Fprintf(ts.w, "%s\t%s\t%s\t%d\t%d\t%d\t\n", f.Id, ep.Region, f.Name, f.Ram, f.Disk, f.VCpus)
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

func findEndpointForRegion(eps []identity.Endpoint, r string) (string, error) {
	for _, ep := range eps {
		if ep.Region == r {
			return ep.PublicURL, nil
		}
	}
	return "", fmt.Errorf("Unknown region %s", r)
}

func TestCreateDestroyServer(t *testing.T) {
	ts, err := setupForList()
	if err != nil {
		t.Error(err)
		return
	}

	imageId := os.Getenv("OS_IMAGE_ID")
	if imageId == "" {
		t.Error("Expected OS_IMAGE_ID environment variable to be set")
		return
	}

	flavorId := os.Getenv("OS_FLAVOR_ID")
	if flavorId == "" {
		t.Error("Expected OS_FLAVOR_ID environment variable to be set")
		return
	}

	region := os.Getenv("OS_REGION_NAME")
	if region == "" {
		region = ts.eps[0].Region
	}

	ep, err := findEndpointForRegion(ts.eps, region)
	if err != nil {
		t.Error(err)
		return
	}

	serverName := randomString("ACPTTEST", 16)
	fmt.Printf("Attempting to create server: %s\n", serverName)

	client := servers.NewClient(ep, ts.a, ts.o)

	cr, err := servers.Create(client, map[string]interface{}{
		"flavorRef": flavorId,
		"imageRef": imageId,
		"name": serverName,
	})
	if err != nil {
		t.Error(err)
		return
	}

	createdServer, err := servers.GetServer(cr)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		servers.Delete(client, createdServer.Id)
	}()

	timeout := 300
	for ; timeout > 0; timeout-- {
		gr, err := servers.GetDetail(client, createdServer.Id)
		if err != nil {
			t.Error(err)
			return
		}

		gottenServer, err := servers.GetServer(gr)
		if err != nil {
			t.Error(err)
			return
		}

		if gottenServer.Id != createdServer.Id {
			t.Error("Created server ID (%s) != gotten server ID (%s)", createdServer.Id, gottenServer.Id)
			return
		}

		if gottenServer.Status == "ACTIVE" {
			fmt.Printf("Server created after %d seconds (approximately)\n", 300-timeout)
			break
		}
		time.Sleep(1*time.Second)
	}
	if timeout < 1 {
		fmt.Printf("I'm not waiting around.\n")
	}
}

func TestUpdateServer(t *testing.T) {
	ts, err := setupForList()
	if err != nil {
		t.Error(err)
		return
	}

	imageId := os.Getenv("OS_IMAGE_ID")
	if imageId == "" {
		t.Error("Expected OS_IMAGE_ID environment variable to be set")
		return
	}

	flavorId := os.Getenv("OS_FLAVOR_ID")
	if flavorId == "" {
		t.Error("Expected OS_FLAVOR_ID environment variable to be set")
		return
	}

	region := os.Getenv("OS_REGION_NAME")
	if region == "" {
		region = ts.eps[0].Region
	}

	ep, err := findEndpointForRegion(ts.eps, region)
	if err != nil {
		t.Error(err)
		return
	}

	serverName := randomString("ACPTTEST", 16)
	fmt.Printf("Attempting to create server: %s\n", serverName)

	client := servers.NewClient(ep, ts.a, ts.o)

	cr, err := servers.Create(client, map[string]interface{}{
		"flavorRef": flavorId,
		"imageRef": imageId,
		"name": serverName,
	})
	if err != nil {
		t.Error(err)
		return
	}

	createdServer, err := servers.GetServer(cr)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		servers.Delete(client, createdServer.Id)
	}()

	timeout := 300
	for ; timeout > 0; timeout-- {
		gr, err := servers.GetDetail(client, createdServer.Id)
		if err != nil {
			t.Error(err)
			return
		}

		gottenServer, err := servers.GetServer(gr)
		if err != nil {
			t.Error(err)
			return
		}

		if gottenServer.Id != createdServer.Id {
			t.Error("Created server ID (%s) != gotten server ID (%s)", createdServer.Id, gottenServer.Id)
			return
		}

		if gottenServer.Status == "ACTIVE" {
			fmt.Printf("Server created after %d seconds (approximately)\n", 300-timeout)
			break
		}
		time.Sleep(1*time.Second)
	}
	if timeout < 1 {
		fmt.Printf("I'm not waiting around.\n")
	}

	alternateName := randomString("ACPTTEST", 16)
	for alternateName == serverName {
		alternateName = randomString("ACPTTEST", 16)
	}

	fmt.Println("Attempting to change server name")

	ur, err := servers.Update(client, createdServer.Id, map[string]interface{}{
		"name": alternateName,
	})
	if err != nil {
		t.Error(err)
		return
	}

	updatedServer, err := servers.GetServer(ur)
	if err != nil {
		t.Error(err)
		return
	}

	if updatedServer.Id != createdServer.Id {
		t.Error("Expected updated and created server to share the same ID")
		return
	}

	timeout = 300
	for ; timeout > 0; timeout-- {
		gr, err := servers.GetDetail(client, createdServer.Id)
		if err != nil {
			t.Error(err)
			return
		}

		gottenServer, err := servers.GetServer(gr)
		if err != nil {
			t.Error(err)
			return
		}

		if gottenServer.Id != updatedServer.Id {
			t.Error("Updated server ID (%s) != gotten server ID (%s)", updatedServer.Id, gottenServer.Id)
			return
		}

		if gottenServer.Name == alternateName {
			fmt.Printf("Server updated after %d seconds (approximately)\n", 300-timeout)
			break
		}
		time.Sleep(1*time.Second)
	}
	if timeout < 1 {
		fmt.Printf("I'm not waiting around.\n")
	}
}

// randomString generates a string of given length, but random content.
// All content will be within the ASCII graphic character set.
// (Implementation from Even Shaw's contribution on
// http://stackoverflow.com/questions/12771930/what-is-the-fastest-way-to-generate-a-long-random-string-in-go).
func randomString(prefix string, n int) string {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return prefix + string(bytes)
}

