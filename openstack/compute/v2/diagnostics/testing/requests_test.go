package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/diagnostics"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetDiagnostics(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDiagnosticGetSuccessfully(t, fakeServer)

	expected := map[string]any{"cpu0_time": float64(173), "memory": float64(524288)}

	res, err := diagnostics.Get(context.TODO(), client.ServiceClient(fakeServer), "1234asdf").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, expected, res)
}
