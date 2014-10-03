// +build acceptance

package v1

import (
	"testing"
	"time"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumetypes"
	"github.com/rackspace/gophercloud/pagination"
)

func TestVolumeTypes(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Failed to create Block Storage v1 client: %v", err)
	}

	vt, err := volumetypes.Create(client, &volumetypes.CreateOpts{
		ExtraSpecs: map[string]interface{}{
			"capabilities": "gpu",
			"priority":     3,
		},
		Name: "gophercloud-test-volumeType",
	}).Extract()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		time.Sleep(10000 * time.Millisecond)
		err = volumetypes.Delete(client, vt.ID)
		if err != nil {
			t.Error(err)
			return
		}
	}()
	t.Logf("Created volume type: %+v\n", vt)

	vt, err = volumetypes.Get(client, vt.ID).Extract()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Got volume type: %+v\n", vt)

	err = volumetypes.List(client).EachPage(func(page pagination.Page) (bool, error) {
		volTypes, err := volumetypes.ExtractVolumeTypes(page)
		if len(volTypes) != 1 {
			t.Errorf("Expected 1 volume type, got %d", len(volTypes))
		}
		t.Logf("Listing volume types: %+v\n", volTypes)
		return true, err
	})
	if err != nil {
		t.Errorf("Error trying to list volume types: %v", err)
	}
}
