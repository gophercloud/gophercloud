// +build acceptance

package v1

import (
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

func TestSnapshots(t *testing.T) {

	client, err := newClient()
	if err != nil {
		t.Fatalf("Failed to create Block Storage v1 client: %v", err)
	}

	v, err := volumes.Create(client, &volumes.CreateOpts{
		Name: "gophercloud-test-volume",
		Size: 1,
	}).Extract()
	if err != nil {
		t.Fatalf("Failed to create volume: %v\n", err)
	}

	err = volumes.WaitForStatus(client, v.ID, "available", 120)
	if err != nil {
		t.Fatalf("Failed to create volume: %v\n", err)
	}

	t.Logf("Created volume: %v\n", v)

	ss, err := snapshots.Create(client, &snapshots.CreateOpts{
		Name:     "gophercloud-test-snapshot",
		VolumeID: v.ID,
	}).Extract()
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v\n", err)
	}

	err = snapshots.WaitForStatus(client, ss.ID, "available", 120)
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v\n", err)
	}

	t.Logf("Created snapshot: %+v\n", ss)

	res = snapshots.Delete(client, ss.ID)
	if res.Err != nil {
		t.Fatalf("Failed to delete snapshot: %v", err)
	}

	err = gophercloud.WaitFor(120, func() (bool, error) {
		_, err := snapshots.Get(client, ss.ID).Extract()
		if err != nil {
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		t.Fatalf("Failed to delete snapshot: %v", err)
	}

	t.Log("Deleted snapshot\n")

	res = volumes.Delete(client, v.ID)
	if res.Err != nil {
		t.Errorf("Failed to delete volume: %v", err)
	}

	err = gophercloud.WaitFor(120, func() (bool, error) {
		_, err := volumes.Get(client, v.ID).Extract()
		if err != nil {
			return true, nil
		}

		return false, nil
	})
	if err != nil {
		t.Errorf("Failed to delete volume: %v", err)
	}

	t.Log("Deleted volume\n")
}
