package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/extensions/injectnetworkinfo"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const serverID = "b16ba811-199d-4ffd-8839-ba96c1185a67"

func TestInjectNetworkInfo(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockInjectNetworkInfoResponse(t, serverID)

	err := injectnetworkinfo.InjectNetworkInfo(context.TODO(), client.ServiceClient(), serverID).ExtractErr()
	th.AssertNoErr(t, err)
}
