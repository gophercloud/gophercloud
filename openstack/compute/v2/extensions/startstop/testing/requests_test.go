package testing

import (
	"testing"

	"github.com/bizflycloud/gophercloud/openstack/compute/v2/extensions/startstop"
	th "github.com/bizflycloud/gophercloud/testhelper"
	"github.com/bizflycloud/gophercloud/testhelper/client"
)

const serverID = "{serverId}"

func TestStart(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockStartServerResponse(t, serverID)

	err := startstop.Start(client.ServiceClient(), serverID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestStop(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockStopServerResponse(t, serverID)

	err := startstop.Stop(client.ServiceClient(), serverID).ExtractErr()
	th.AssertNoErr(t, err)
}
