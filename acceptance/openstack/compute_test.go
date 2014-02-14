package openstack

import (
	"fmt"
	"github.com/rackspace/gophercloud/openstack/compute/flavors"
	"github.com/rackspace/gophercloud/openstack/compute/images"
	"github.com/rackspace/gophercloud/openstack/compute/servers"
	"os"
	"testing"
)

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

func TestCreateDestroyServer(t *testing.T) {
	ts, err := setupForCRUD()
	if err != nil {
		t.Error(err)
		return
	}

	err = createServer(ts)
	if err != nil {
		t.Error(err)
		return
	}

	// We put this in a defer so that it gets executed even in the face of errors or panics.
	defer func() {
		servers.Delete(ts.client, ts.createdServer.Id)
	}()

	err = waitForStatus(ts, "ACTIVE")
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateServer(t *testing.T) {
	ts, err := setupForCRUD()
	if err != nil {
		t.Error(err)
		return
	}

	err = createServer(ts)
	if err != nil {
		t.Error(err)
		return
	}

	defer func() {
		servers.Delete(ts.client, ts.createdServer.Id)
	}()

	err = waitForStatus(ts, "ACTIVE")
	if err != nil {
		t.Error(err)
		return
	}

	err = changeServerName(ts)
	if err != nil {
		t.Error(err)
		return
	}
}
