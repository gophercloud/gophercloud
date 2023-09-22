//go:build acceptance || blockstorage
// +build acceptance blockstorage

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumetypes"
)

func TestVolumeTypesList(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/icehouse")
	client, err := clients.NewBlockStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create a blockstorage client: %v", err)
	}

	allPages, err := volumetypes.List(client).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve volume types: %v", err)
	}

	allVolumeTypes, err := volumetypes.ExtractVolumeTypes(allPages)
	if err != nil {
		t.Fatalf("Unable to extract volume types: %v", err)
	}

	for _, volumeType := range allVolumeTypes {
		tools.PrintResource(t, volumeType)
	}
}

func TestVolumeTypesCreateDestroy(t *testing.T) {
	clients.SkipReleasesAbove(t, "stable/icehouse")
	client, err := clients.NewBlockStorageV1Client()
	if err != nil {
		t.Fatalf("Unable to create a blockstorage client: %v", err)
	}

	volumeType, err := CreateVolumeType(t, client)
	if err != nil {
		t.Fatalf("Unable to create volume type: %v", err)
	}
	defer DeleteVolumeType(t, client, volumeType)

	tools.PrintResource(t, volumeType)
}
