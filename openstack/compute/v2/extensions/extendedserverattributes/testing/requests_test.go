package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/extendedserverattributes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestServerWithUsageExt(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/servers/d650a0ce-17c3-497d-961a-43c4af80998a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, ServerWithAttributesExtResult)
	})

	type serverAttributesExt struct {
		servers.Server
		extendedserverattributes.ServerAttributesExt
	}
	var serverWithAttributesExt serverAttributesExt
	err := servers.Get(fake.ServiceClient(), "d650a0ce-17c3-497d-961a-43c4af80998a").ExtractInto(&serverWithAttributesExt)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, serverWithAttributesExt.ReservationID, "r-ky9gim1l")
	th.AssertEquals(t, serverWithAttributesExt.LaunchIndex, 0)
	th.AssertEquals(t, serverWithAttributesExt.Hostname, "test00")
	th.AssertEquals(t, serverWithAttributesExt.KernelID, "")
	th.AssertEquals(t, serverWithAttributesExt.RamdiskID, "")
	th.AssertEquals(t, serverWithAttributesExt.Host, "compute01")
	th.AssertEquals(t, serverWithAttributesExt.RootDeviceName, "/dev/sda")
	th.AssertEquals(t, serverWithAttributesExt.UserData, "")
	th.AssertEquals(t, serverWithAttributesExt.InstanceName, "instance-00000001")
	th.AssertEquals(t, serverWithAttributesExt.HypervisorHostname, "compute01")
	th.AssertEquals(t, serverWithAttributesExt.Created, time.Date(2018, 07, 27, 9, 15, 48, 0, time.UTC))
	th.AssertEquals(t, serverWithAttributesExt.Updated, time.Date(2018, 07, 27, 9, 15, 55, 0, time.UTC))
	th.AssertEquals(t, serverWithAttributesExt.ID, "d650a0ce-17c3-497d-961a-43c4af80998a")
	th.AssertEquals(t, serverWithAttributesExt.Name, "test_instance")
	th.AssertEquals(t, serverWithAttributesExt.Status, "ACTIVE")
	th.AssertEquals(t, serverWithAttributesExt.UserID, "0f2f3822679e4b3ea073e5d1c6ed5f02")
	th.AssertEquals(t, serverWithAttributesExt.TenantID, "424e7cf0243c468ca61732ba45973b3e")
}
