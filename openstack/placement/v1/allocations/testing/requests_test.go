package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/allocations"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetAllocationsSuccess(t, fakeServer)

	actual, err := allocations.Get(context.TODO(), client.ServiceClient(fakeServer), ConsumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocations, *actual)
}

func TestGetEmptySuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetEmptyAllocationsSuccess(t, fakeServer)

	actual, err := allocations.Get(context.TODO(), client.ServiceClient(fakeServer), EmptyConsumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedEmptyAllocations, *actual)
}
