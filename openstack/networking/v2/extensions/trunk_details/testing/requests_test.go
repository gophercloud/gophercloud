package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunk_details"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestServerWithUsageExt(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	const portIDFixture = "dc3e8758-ee96-402d-94b0-4be5e9396c82"

	fakeServer.Mux.HandleFunc("/ports/"+portIDFixture, func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprint(w, PortWithTrunkDetailsResult)
	})

	var portExt struct {
		ports.Port
		trunk_details.TrunkDetailsExt
	}

	// Extract basic fields.
	err := ports.Get(context.TODO(), client.ServiceClient(fakeServer), portIDFixture).ExtractInto(&portExt)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, portExt.TrunkID, "f170c831-8c55-4ceb-ad13-75eab4a121e5")
	th.AssertEquals(t, len(portExt.SubPorts), 1)
	subPort := portExt.SubPorts[0]
	th.AssertEquals(t, subPort.SegmentationID, 100)
	th.AssertEquals(t, subPort.SegmentationType, "vlan")
	th.AssertEquals(t, subPort.PortID, "20c673d8-7f9d-4570-b662-148d9ddcc5bd")
	th.AssertEquals(t, subPort.MACAddress, "fa:16:3e:88:29:a0")
}
