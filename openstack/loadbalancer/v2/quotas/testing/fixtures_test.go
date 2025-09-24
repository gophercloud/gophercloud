package testing

import (
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/quotas"
	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

const GetResponseRaw_1 = `
{
    "quota": {
        "loadbalancer": 15,
        "listener": 30,
        "member": -1,
        "pool": 15,
        "healthmonitor": 30,
        "l7policy": 100,
        "l7rule": -1
    }
}
`

const GetResponseRaw_2 = `
{
    "quota": {
        "load_balancer": 15,
        "listener": 30,
        "member": -1,
        "pool": 15,
        "health_monitor": 30,
        "l7policy": 100,
        "l7rule": -1
    }
}
`

var GetResponse = quotas.Quota{
	Loadbalancer:  15,
	Listener:      30,
	Member:        -1,
	Pool:          15,
	Healthmonitor: 30,
	L7Policy:      100,
	L7Rule:        -1,
}

const UpdateRequestResponseRaw_1 = `
{
    "quota": {
        "loadbalancer": 20,
        "listener": 40,
        "member": 200,
        "pool": 20,
        "healthmonitor": -1,
        "l7policy": 50,
        "l7rule": 100
    }
}
`

const UpdateRequestResponseRaw_2 = `
{
    "quota": {
        "load_balancer": 20,
        "listener": 40,
        "member": 200,
        "pool": 20,
        "health_monitor": -1,
        "l7policy": 50,
        "l7rule": 100
    }
}
`

var UpdateResponse = quotas.Quota{
	Loadbalancer:  20,
	Listener:      40,
	Member:        200,
	Pool:          20,
	Healthmonitor: -1,
	L7Policy:      50,
	L7Rule:        100,
}

func MockDeleteResponse(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/v2.0/quotas/0a73845280574ad389c292f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}
