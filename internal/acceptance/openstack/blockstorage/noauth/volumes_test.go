//go:build acceptance || blockstorage || volumes

package noauth

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
)

func TestVolumesList(t *testing.T) {
	clients.RequireCinderNoAuth(t)

	client, err := clients.NewBlockStorageV3NoAuthClient()
	if err != nil {
		t.Fatalf("Unable to create a blockstorage client: %v", err)
	}

	allPages, err := volumes.List(client, volumes.ListOpts{}).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to retrieve volumes: %v", err)
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		t.Fatalf("Unable to extract volumes: %v", err)
	}

	for _, volume := range allVolumes {
		tools.PrintResource(t, volume)
	}
}

func TestVolumesCreateDestroy(t *testing.T) {
	clients.RequireCinderNoAuth(t)

	client, err := clients.NewBlockStorageV3NoAuthClient()
	if err != nil {
		t.Fatalf("Unable to create blockstorage client: %v", err)
	}

	volume, err := CreateVolume(t, client)
	if err != nil {
		t.Fatalf("Unable to create volume: %v", err)
	}
	defer DeleteVolume(t, client, volume)

	newVolume, err := volumes.Get(context.TODO(), client, volume.ID).Extract()
	if err != nil {
		t.Errorf("Unable to retrieve volume: %v", err)
	}

	tools.PrintResource(t, newVolume)
}
