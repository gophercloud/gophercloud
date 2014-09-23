package quotas

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	th "github.com/rackspace/gophercloud/testhelper"
)

const tokenID = "123"

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenID},
		Endpoint: th.Endpoint(),
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/quotas", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "quota": {
        "subnet": 10,
        "router": 10,
        "port": 50,
        "network": 10,
        "floatingip": 50
    }
}
      `)
	})

	qs, err := Get(serviceClient()).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, qs.Subnet, 10)
	th.AssertEquals(t, qs.Router, 10)
	th.AssertEquals(t, qs.Port, 50)
	th.AssertEquals(t, qs.Network, 10)
	th.AssertEquals(t, qs.FloatingIP, 50)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/quotas", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "quota": {
        "subnet": 40,
        "router": 40,
        "network": 10,
        "floatingip": 30,
        "port": 30
    }
}
			`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
    "quota": {
        "subnet": 40,
        "router": 40,
        "network": 10,
        "floatingip": 30,
        "port": 30
    }
}
    `)
	})

	i10, i30, i40 := 10, 30, 40
	opts := UpdateOpts{Subnet: &i40, Router: &i40, Network: &i10, FloatingIP: &i30, Port: &i30}
	qs, err := Update(serviceClient(), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, qs.Subnet, 40)
	th.AssertEquals(t, qs.Router, 40)
	th.AssertEquals(t, qs.Port, 30)
	th.AssertEquals(t, qs.Network, 10)
	th.AssertEquals(t, qs.FloatingIP, 30)
}

func TestReset(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/quotas", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := Reset(serviceClient())
	th.AssertNoErr(t, res.Err)
}
