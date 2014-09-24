// +build acceptance

package v1

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
)

func waitForSnapshot(client *gophercloud.ServiceClient, id string) error {
	for secondsSlept := 0; secondsSlept < 240; secondsSlept++ {
		fmt.Printf("Seconds slept waiting for snapshot: %d\n", secondsSlept)
		gss, err := snapshots.Get(client, id).ExtractSnapshot()
		if err != nil {
			return err
		}
		if gss.Status == "available" {
			return nil
		}
		if gss.Status == "error" {
			return fmt.Errorf("Error waiting for snapshot to create. Snapshot status is 'error'.")
		}
		time.Sleep(1 * time.Second)
	}
	gss, err := snapshots.Get(client, id).ExtractSnapshot()
	if err != nil {
		return err
	}
	return fmt.Errorf("Time out waiting for snapshot to become available: %+v", gss)
}

func TestSnapshots(t *testing.T) {

	volumeID := os.Getenv("OS_VOLUME_ID")
	if volumeID == "" {
		t.Errorf("Expect OS_VOLUME_ID environment variable. Skipping create and delete snapshot functions.")
		return
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Failed to create Block Storage v1 client: %v", err)
	}

	css, err := snapshots.Create(client, snapshots.CreateOpts{
		Name:     "gophercloud-test-snapshot",
		VolumeID: volumeID,
	})
	if err != nil {
		t.Errorf("Failed to create snapshot: %v\n", err)
	}

	err = waitForSnapshot(client, css.ID)
	if err != nil {
		t.Errorf("Failed to create snapshot: %v\n", err)
	}

	err = waitForSnapshot(client, css.ID)

	t.Logf("Created snapshots: %+v\n", *css)

}
