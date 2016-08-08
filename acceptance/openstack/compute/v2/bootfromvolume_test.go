// +build acceptance compute bootfromvolume

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
)

func TestBootFromVolumeSingleVolume(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := clients.AcceptanceTestChoicesFromEnv()
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

	server, err := CreateBootableVolumeServer(t, client, blockDevices, choices)
	if err != nil {
		t.Fatal("Unable to create server: %v", err)
	}
	defer DeleteServer(t, client, server)

	PrintServer(t, server)
}

func TestBootFromVolumeMultiEphemeral(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires server creation in short mode.")
	}

	client, err := clients.NewComputeV2Client()
	if err != nil {
		t.Fatalf("Unable to create a compute client: %v", err)
	}

	choices, err := clients.AcceptanceTestChoicesFromEnv()
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

	server, err := CreateBootableVolumeServer(t, client, blockDevices, choices)
	if err != nil {
		t.Fatalf("Unable to create server: %v", err)
	}
	defer DeleteServer(t, client, server)

	PrintServer(t, server)
}
