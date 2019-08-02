package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/portforwarding"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/floatingips/2f95fd2b-9f6a-4e8e-9e9a-2cbe286cbf9e/port_forwardings", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
  	"port_forwarding": {
      "protocol": "tcp",
      "internal_ip_address": "10.0.0.11",
      "internal_port": 25,
      "internal_port_id": "1238be08-a2a8-4b8d-addf-fb5e2250e480",
      "external_port": 2230
  }
}
		
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, `
{
	"port_forwarding": {
    		"protocol": "tcp",
    		"internal_ip_address": "10.0.0.11",
    		"internal_port": 25,
    		"internal_port_id": "1238be08-a2a8-4b8d-addf-fb5e2250e480",
    		"external_port": 2230,
    		"id": "725ade3c-9760-4880-8080-8fc2dbab9acc"
  }
}`)
	})

	options := portforwarding.CreateOpts{
		Protocol:          "tcp",
		InternalIPAddress: "10.0.0.11",
		InternalPort:      25,
		ExternalPort:      2230,
		InternalPortID:    "1238be08-a2a8-4b8d-addf-fb5e2250e480",
	}

	pf, err := portforwarding.Create(fake.ServiceClient(), "2f95fd2b-9f6a-4e8e-9e9a-2cbe286cbf9e", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "725ade3c-9760-4880-8080-8fc2dbab9acc", pf.ID)
	th.AssertEquals(t, "10.0.0.11", pf.InternalIPAddress)
	th.AssertEquals(t, 25, pf.InternalPort)
	th.AssertEquals(t, "1238be08-a2a8-4b8d-addf-fb5e2250e480", pf.InternalPortID)
	th.AssertEquals(t, 2230, pf.ExternalPort)
	th.AssertEquals(t, "tcp", pf.Protocol)
}
