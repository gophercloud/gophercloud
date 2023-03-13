//go:build acceptance || blockstorage
// +build acceptance blockstorage

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/availabilityzones"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestAvailabilityZonesList(t *testing.T) {
	clients.RequireLong(t)

	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	allPages, err := availabilityzones.List(client).AllPages()
	th.AssertNoErr(t, err)

	zones, err := availabilityzones.ExtractAvailabilityZones(allPages)
	th.AssertNoErr(t, err)

	if len(zones) == 0 {
		t.Fatal("At least one availability zone was expected to be found")
	}
}
