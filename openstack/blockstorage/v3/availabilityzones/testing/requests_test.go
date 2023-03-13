package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/availabilityzones"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// Verifies that availability zones can be listed correctly
func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	allPages, err := availabilityzones.List(client.ServiceClient()).AllPages()
	th.AssertNoErr(t, err)
	actual, err := availabilityzones.ExtractAvailabilityZones(allPages)
	th.AssertNoErr(t, err)
	expected := []availabilityzones.AvailabilityZone{
		{
			Name:  "nova",
			State: availabilityzones.AvailabilityZoneState{Available: true},
		},
	}

	th.CheckDeepEquals(t, expected, actual)
}
