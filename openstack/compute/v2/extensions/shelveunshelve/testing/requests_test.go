package testing

import (
	"testing"

	"github.com/bizflycloud/gophercloud/openstack/compute/v2/extensions/shelveunshelve"
	th "github.com/bizflycloud/gophercloud/testhelper"
	"github.com/bizflycloud/gophercloud/testhelper/client"
)

const serverID = "{serverId}"
const availabilityZone = "test-zone"

func TestShelve(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockShelveServerResponse(t, serverID)

	err := shelveunshelve.Shelve(client.ServiceClient(), serverID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestShelveOffload(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockShelveOffloadServerResponse(t, serverID)

	err := shelveunshelve.ShelveOffload(client.ServiceClient(), serverID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUnshelveNoAvailabilityZone(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	unshelveOpts := shelveunshelve.UnshelveOpts{}

	mockUnshelveServerResponseNoAvailabilityZone(t, serverID)

	err := shelveunshelve.Unshelve(client.ServiceClient(), serverID, unshelveOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUnshelveWithAvailabilityZone(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	unshelveOpts := shelveunshelve.UnshelveOpts{
		AvailabilityZone: availabilityZone,
	}

	mockUnshelveServerResponseWithAvailabilityZone(t, serverID, availabilityZone)

	err := shelveunshelve.Unshelve(client.ServiceClient(), serverID, unshelveOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
