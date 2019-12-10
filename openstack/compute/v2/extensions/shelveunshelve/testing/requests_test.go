package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/shelveunshelve"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const serverID = "{serverId}"

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

func TestUnshelve(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockUnshelveServerResponse(t, serverID)

	err := shelveunshelve.Unshelve(client.ServiceClient(), serverID).ExtractErr()
	th.AssertNoErr(t, err)
}
