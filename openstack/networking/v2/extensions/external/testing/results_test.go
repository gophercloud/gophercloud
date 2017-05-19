package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/networks/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "network": {
        "admin_state_up": true,
        "id": "8d05a1b1-297a-46ca-8974-17debf51ca3c",
        "name": "ext_net",
        "router:external": true,
        "shared": false,
        "status": "ACTIVE",
        "subnets": [
            "2f1fb918-9b0e-4bf9-9a50-6cebbb4db2c5"
        ],
        "tenant_id": "5eb8995cf717462c9df8d1edfa498010"
    }
}
			`)
	})

	type NetworkPlusExternal struct {
		networks.Network
		external.NetworkExt
	}

	var actual NetworkPlusExternal

	err := networks.Get(fake.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").ExtractInto(&actual)

	th.AssertNoErr(t, err)
	th.AssertEquals(t, true, actual.External)
}
