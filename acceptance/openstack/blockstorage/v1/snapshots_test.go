// +build acceptance

package v1

import (
	"strconv"
	"testing"
	//"time"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

func waitForVolume(id string) {

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

	waitForVolume(cv.ID)

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
