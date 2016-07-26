// +build acceptance compute bootfromvolume

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func TestBootFromVolumeSingleVolume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	blockDevices := []bootfromvolume.BlockDevice{
		bootfromvolume.BlockDevice{
			UUID:       choices.ImageID,
			SourceType: bootfromvolume.Image,
			VolumeSize: 10,
		},
	}

	server, err := createBootableVolumeServer(t, client, blockDevices, choices)
	if err != nil {
		t.Fatal("Unable to create server: %v", err)
	}
	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatal("Unable to wait for server: %v", err)
	}
	defer deleteServer(t, client, server)

	printServer(t, server)
}

func TestBootFromVolumeMultiEphemeral(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	blockDevices := []bootfromvolume.BlockDevice{
		bootfromvolume.BlockDevice{
			BootIndex:           0,
			UUID:                choices.ImageID,
			SourceType:          bootfromvolume.Image,
			DestinationType:     "local",
			DeleteOnTermination: true,
			VolumeSize:          5,
		},
		bootfromvolume.BlockDevice{
			BootIndex:           -1,
			SourceType:          bootfromvolume.Blank,
			DestinationType:     "local",
			DeleteOnTermination: true,
			GuestFormat:         "ext4",
			VolumeSize:          1,
		},
		bootfromvolume.BlockDevice{
			BootIndex:           -1,
			SourceType:          bootfromvolume.Blank,
			DestinationType:     "local",
			DeleteOnTermination: true,
			GuestFormat:         "ext4",
			VolumeSize:          1,
		},
	}

	server, err := createBootableVolumeServer(t, client, blockDevices, choices)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}
	if err = waitForStatus(client, server, "ACTIVE"); err != nil {
		t.Fatalf("Unable to wait for server: %v", err)
	}
	defer deleteServer(t, client, server)

	printServer(t, server)
}

func createBootableVolumeServer(t *testing.T, client *gophercloud.ServiceClient, blockDevices []bootfromvolume.BlockDevice, choices *ComputeChoices) (*servers.Server, error) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	networkID, err := getNetworkIDFromTenantNetworks(t, client, choices.NetworkName)
	if err != nil {
		t.Fatalf("Failed to obtain network ID: %v", err)
	}

	name := tools.RandomString("ACPTTEST", 16)
	t.Logf("Attempting to create bootable volume server: %s", name)

	serverCreateOpts := servers.CreateOpts{
		Name:      name,
		FlavorRef: choices.FlavorID,
		ImageRef:  choices.ImageID,
		Networks: []servers.Network{
			servers.Network{UUID: networkID},
		},
	}

	server, err := bootfromvolume.Create(client, bootfromvolume.CreateOptsExt{
		serverCreateOpts,
		blockDevices,
	}).Extract()

	if err != nil {
		return server, err
	}

	return server, nil
}
