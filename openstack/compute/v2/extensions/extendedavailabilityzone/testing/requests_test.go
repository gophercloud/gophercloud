package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/extendedavailabilityzone"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestServerServerWithAvailabilityZoneExt(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/d650a0ce-17c3-497d-961a-43c4af80998a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, ServerWithAvailabilityZoneExtResult)
	})

	type serverAvailabilityZoneExt struct {
		servers.Server
		extendedavailabilityzone.AvailabilityZoneExt
	}
	var serverWithAvailabilityZoneExt serverAvailabilityZoneExt
	err := servers.Get(fake.ServiceClient(), "d650a0ce-17c3-497d-961a-43c4af80998a").ExtractInto(&serverWithAvailabilityZoneExt)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, serverWithAvailabilityZoneExt.AvailabilityZone, "nova")
	th.AssertEquals(t, serverWithAvailabilityZoneExt.Created, time.Date(2018, 07, 27, 9, 15, 48, 0, time.UTC))
	th.AssertEquals(t, serverWithAvailabilityZoneExt.Updated, time.Date(2018, 07, 27, 9, 15, 55, 0, time.UTC))
	th.AssertEquals(t, serverWithAvailabilityZoneExt.ID, "d650a0ce-17c3-497d-961a-43c4af80998a")
	th.AssertEquals(t, serverWithAvailabilityZoneExt.Name, "test_instance")
	th.AssertEquals(t, serverWithAvailabilityZoneExt.Status, "ACTIVE")
	th.AssertEquals(t, serverWithAvailabilityZoneExt.UserID, "0f2f3822679e4b3ea073e5d1c6ed5f02")
	th.AssertEquals(t, serverWithAvailabilityZoneExt.TenantID, "424e7cf0243c468ca61732ba45973b3e")
}
