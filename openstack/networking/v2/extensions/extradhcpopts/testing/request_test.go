package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestGetWithDHCPOpts(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/ports/46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, GetWithDHCPOptsResponse)
	})

	var s struct {
		ports.Port
		extradhcpopts.DHCPOptsExt
	}

	err := ports.Get(fake.ServiceClient(), "46d4bfb9-b26e-41f3-bd2e-e6dcc1ccedb2").ExtractInto(&s)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, s.Status, "ACTIVE")
	th.AssertEquals(t, s.NetworkID, "a87cc70a-3e15-4acf-8205-9b711a3531b7")
	th.AssertEquals(t, s.TenantID, "d6700c0c9ffa4f1cb322cd4a1f3906fa")
	th.AssertDeepEquals(t, s.DHCPOptsExt, extradhcpopts.DHCPOptsExt{
		DHCPOpts: []extradhcpopts.DHCPOpts{
			{DHCPOptName: "option1", DHCPOptValue: "value1", DHCPOptIPVersion: 4},
			{DHCPOptName: "option2", DHCPOptValue: "value2", DHCPOptIPVersion: 4},
		},
	})
	th.AssertEquals(t, s.AdminStateUp, true)
	th.AssertEquals(t, s.Name, "port-with-extra-dhcp-opts")
	th.AssertEquals(t, s.DeviceOwner, "")
	th.AssertEquals(t, s.MACAddress, "fa:16:3e:c9:cb:f0")
	th.AssertDeepEquals(t, s.FixedIPs, []ports.IP{
		{SubnetID: "a0304c3a-4f08-4c43-88af-d796509c97d2", IPAddress: "10.0.0.4"},
	})
	th.AssertEquals(t, s.ID, "65c0ee9f-d634-4522-8954-51021b570b0d")
	th.AssertEquals(t, s.DeviceID, "")
}
