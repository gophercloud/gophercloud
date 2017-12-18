// +build acceptance blockstorage

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
)

func TestVolumeTypesList(t *testing.T) {
	client, err := clients.NewBlockStorageV3Client()
	if err != nil {
		t.Fatalf("Unable to create a blockstorage client: %v", err)
	}

	listOpts := volumetypes.ListOpts{
		Sort:  "name:asc",
		Limit: 1,
	}

	allPages, err := volumetypes.List(client, listOpts).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve volumetypes: %v", err)
	}

	allVolumes, err := volumetypes.ExtractVolumeTypes(allPages)
	if err != nil {
		t.Fatalf("Unable to extract volumetypes: %v", err)
	}

	for _, volume := range allVolumes {
		tools.PrintResource(t, volume)
	}

	if len(allVolumes) > 0 {
		vt, err := volumetypes.Get(client, allVolumes[0].ID).Extract()
		if err != nil {
			t.Fatalf("Error retrieving volume type: %v", err)
		}

		tools.PrintResource(t, vt)
	}
}
