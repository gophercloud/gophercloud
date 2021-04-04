package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/qos"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

var createQoSExpected = qos.QoS{
	ID:       "d32019d3-bc6e-4319-9c1d-6722fc136a22",
	Name:     "qos-001",
	Consumer: "front-end",
	Specs: map[string]string{
		"read_iops_sec": "20000",
	},
}

func MockCreateResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `
{
    "qos_specs": {
    	"name": "qos-001",
		"consumer": "front-end",
		"read_iops_sec": "20000"
    }
}
      `)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, `
{
  "qos_specs": {
    "id": "d32019d3-bc6e-4319-9c1d-6722fc136a22",
	"name": "qos-001",
	"consumer": "front-end",
	"specs": {
	  "read_iops_sec": "20000"
	}
  }
}
    `)
	})
}

func MockDeleteResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}
