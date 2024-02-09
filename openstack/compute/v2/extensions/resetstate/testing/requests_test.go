package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/resetstate"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"

func TestResetState(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockResetStateResponse(t, serverID, "active")

	err := resetstate.ResetState(context.TODO(), client.ServiceClient(), serverID, "active").ExtractErr()
	th.AssertNoErr(t, err)
}
