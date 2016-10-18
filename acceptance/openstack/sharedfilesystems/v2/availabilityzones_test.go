package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/availabilityzones"
)

func TestAvailabilityZonesList(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	zones, err := availabilityzones.List(client).Extract()
	if err != nil {
		t.Fatalf("Unable to list availability zones: %v", err)
	}

	if len(zones) == 0 {
		t.Fatal("At least one availability zone was expected to be found")
	}
}
