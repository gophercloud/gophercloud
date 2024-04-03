//go:build acceptance || sharedfilesystems || availabilityzones

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/availabilityzones"
)

func TestAvailabilityZonesList(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	if err != nil {
		t.Fatalf("Unable to create shared file system client: %v", err)
	}

	allPages, err := availabilityzones.List(client).AllPages(context.TODO())
	if err != nil {
		t.Fatalf("Unable to list availability zones: %v", err)
	}

	zones, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		t.Fatalf("Unable to extract availability zones: %v", err)
	}

	if len(zones) == 0 {
		t.Fatal("At least one availability zone was expected to be found")
	}
}
