// +build acceptance

package v1

import (
	"testing"
	"time"

	"github.com/rackspace/gophercloud/openstack/blockStorage/v1/volumeTypes"
	"github.com/rackspace/gophercloud/pagination"
)

func TestVolumeTypes(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Failed to create Block Storage v1 client: %v", err)
	}

	vt, err := volumeTypes.Create(client, volumeTypes.CreateOpts{
		ExtraSpecs: map[string]interface{}{
			"capabilities": "gpu",
			"priority":     3,
		},
		Name: "gophercloud-test-volumeType",
	})
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		time.Sleep(10000 * time.Millisecond)
		err = volumeTypes.Delete(client, vt.ID)
		if err != nil {
			t.Error(err)
			return
		}
	}()
	t.Logf("Created volume type: %+v\n", vt)

	vt, err = volumeTypes.Get(client, vt.ID).ExtractVolumeType()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Got volume type: %+v\n", vt)

	err = volumeTypes.List(client, volumeTypes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		volTypes, err := volumeTypes.ExtractVolumeTypes(page)
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
