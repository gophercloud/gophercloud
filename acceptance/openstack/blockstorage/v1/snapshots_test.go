// +build acceptance

package v1

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

func waitForVolume(client *gophercloud.ServiceClient, id string) error {
	notReady := true
	secondsSlept := 0
	for notReady && secondsSlept < 20 {
		gv, err := volumes.Get(client, id).ExtractVolume()
		if err != nil {
			return err
		}
		if gv.Status == "available" {
			return nil
		}
		time.Sleep(1 * time.Millisecond)
		secondsSlept = secondsSlept + 1
	}

	return errors.New("Time out waiting for volume to become available")
}

var numSnapshots = 1

func TestSnapshots(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Failed to create Block Storage v1 client: %v", err)
	}

	cv, err := volumes.Create(client, volumes.CreateOpts{
		Size: 1,
		Name: "gophercloud-test-volume",
	})
	if err != nil {
		t.Fatalf("Failed to create volume: %v", err)
	}

	err = waitForVolume(client, cv.ID)
	if err != nil {
		t.Fatal(err)
	}

	var sss []*snapshots.Snapshot
	for i := 0; i < numSnapshots; i++ {
		css, err := snapshots.Create(client, snapshots.CreateOpts{
			Name:     "gophercloud-test-snapshot-" + strconv.Itoa(i),
			VolumeID: cv.ID,
		})
		if err != nil {
			t.Errorf("Failed to create snapshot: %v\n", err)
		}
		sss = append(sss, css)
	}

	t.Logf("Created snapshots: %+v\n", sss)
}
