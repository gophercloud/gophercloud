package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/resetnetwork"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"

func TestResetNetwork(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockResetNetworkResponse(t, serverID)

	err := resetnetwork.ResetNetwork(context.TODO(), client.ServiceClient(), serverID).ExtractErr()
	th.AssertNoErr(t, err)
}
