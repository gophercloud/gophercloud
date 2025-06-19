package testing

import (
	"context"
	"testing"

	az "github.com/gophercloud/gophercloud/v2/openstack/compute/v2/availabilityzones"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// Verifies that availability zones can be listed correctly
func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetSuccessfully(t, fakeServer)

	allPages, err := az.List(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := az.ExtractAvailabilityZones(allPages)
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, AZResult, actual)
}

// Verifies that detailed availability zones can be listed correctly
func TestListDetail(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetDetailSuccessfully(t, fakeServer)

	allPages, err := az.ListDetail(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := az.ExtractAvailabilityZones(allPages)
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, AZDetailResult, actual)
}
