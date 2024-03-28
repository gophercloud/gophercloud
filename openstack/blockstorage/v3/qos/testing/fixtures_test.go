package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/qos"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

var createQoSExpected = qos.QoS{
	ID:       "d32019d3-bc6e-4319-9c1d-6722fc136a22",
	Name:     "qos-001",
	Consumer: "front-end",
	Specs: map[string]string{
		"read_iops_sec": "20000",
	},
}

var getQoSExpected = qos.QoS{
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

func MockListResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `
					{
						"qos_specs": [
							{
								"consumer": "back-end",
								"id": "1",
								"name": "foo",
								"specs": {}
							},
							{
								"consumer": "front-end",
								"id": "2",
								"name": "bar",
								"specs" : {
									"read_iops_sec" : "20000"
								 }
							}

						],
						"qos_specs_links": [
							{
								"href": "%s/qos-specs?marker=2",
								"rel": "next"
							}
						]
					}
				`, th.Server.URL)
		case "2":
			fmt.Fprintf(w, `{ "qos_specs": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
}

func MockGetResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
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

// UpdateBody provides a PUT result of the qos_specs for a QoS
const UpdateBody = `
{
    "qos_specs" : {
		"consumer": "back-end",
		"read_iops_sec":  "40000",
		"write_iops_sec": "40000"
    }
}
`

// UpdateQos is the expected qos_specs returned from PUT on a qos's qos_specs
var UpdateQos = map[string]string{
	"consumer":       "back-end",
	"read_iops_sec":  "40000",
	"write_iops_sec": "40000",
}

func MockUpdateResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs/d32019d3-bc6e-4319-9c1d-6722fc136a22", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, `{
				"qos_specs": {
					"consumer": "back-end",
					"read_iops_sec":  "40000",
					"write_iops_sec": "40000"
				}
			}`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateBody)
	})
}

func MockDeleteKeysResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs/d32019d3-bc6e-4319-9c1d-6722fc136a22/delete_keys", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `{
			"keys": [
				"read_iops_sec"
			]
		}`)

		w.WriteHeader(http.StatusAccepted)
	})
}

func MockAssociateResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs/d32019d3-bc6e-4319-9c1d-6722fc136a22/associate", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}

func MockDisassociateResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs/d32019d3-bc6e-4319-9c1d-6722fc136a22/disassociate", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}

func MockDisassociateAllResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs/d32019d3-bc6e-4319-9c1d-6722fc136a22/disassociate_all", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusAccepted)
	})
}

func MockListAssociationsResponse(t *testing.T) {
	th.Mux.HandleFunc("/qos-specs/d32019d3-bc6e-4319-9c1d-6722fc136a22/associations", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
			  "qos_associations": [
			    {
			      "name": 			  "foo",
				  "id": 			  "2f954bcf047c4ee9b09a37d49ae6db54",
				  "association_type": "volume_type"
			    }
			  ]
			}
		`)
	})
}
