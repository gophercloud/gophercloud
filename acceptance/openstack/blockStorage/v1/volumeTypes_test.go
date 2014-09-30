// +build acceptance

package v1

import (
	"strconv"
	"testing"
	"time"

	"github.com/rackspace/gophercloud/openstack/blockStorage/v1/volumeTypes"
	"github.com/rackspace/gophercloud/pagination"
)

var numVolTypes = 1

func TestVolumeTypes(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Failed to create Block Storage v1 client: %v", err)
	}

	var cvt *volumeTypes.VolumeType
	for i := 0; i < numVolTypes; i++ {
		cvt, err = volumeTypes.Create(client, volumeTypes.CreateOpts{
			ExtraSpecs: map[string]interface{}{
				"capabilities": "gpu",
				"priority":     3,
			},
			Name: "gophercloud-test-volumeType-" + strconv.Itoa(i),
		})
		if err != nil {
			t.Error(err)
			return
		}
		defer func() {
			time.Sleep(10000 * time.Millisecond)
			err = volumeTypes.Delete(client, cvt.ID)
			if err != nil {
				t.Error(err)
				return
			}
		}()
		t.Logf("Created volume type: %+v\n", cvt)
	}

	gr, err := volumeTypes.Get(client, cvt.ID)
	if err != nil {
		t.Error(err)
		return
	}
	v, err := volumeTypes.ExtractVolumeType(gr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("Got volume type: %+v\n", v)

	pager := volumeTypes.List(client, volumeTypes.ListOpts{})
	if err != nil {
		t.Error(err)
		return
	}
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		volTypes, err := volumeTypes.ExtractVolumeTypes(page)
		if len(volTypes) != numVolTypes {
			t.Errorf("Expected %d volume types, got %d", numVolTypes, len(volTypes))
		}
		t.Logf("Listing volume types: %+v\n", volTypes)
		return true, err
	})

}
