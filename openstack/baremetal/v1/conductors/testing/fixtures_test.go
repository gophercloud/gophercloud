package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/conductors"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ConductorListBody contains the canned body of a conductor.List response, without detail.
const ConductorListBody = `
 {
  "conductors": [
    {
      "hostname": "compute1.localdomain",
      "conductor_group": "",
      "links": [
        {
          "href": "http://127.0.0.1:6385/v1/conductors/compute1.localdomain",
          "rel": "self"
        },
        {
          "href": "http://127.0.0.1:6385/conductors/compute1.localdomain",
          "rel": "bookmark"
        }
      ],
      "alive": false
    },
    {
      "hostname": "compute2.localdomain",
      "conductor_group": "",
      "links": [
        {
          "href": "http://127.0.0.1:6385/v1/conductors/compute2.localdomain",
          "rel": "self"
        },
        {
          "href": "http://127.0.0.1:6385/conductors/compute2.localdomain",
          "rel": "bookmark"
        }
      ],
      "alive": true
    }
  ]
 }
`

// ConductorListDetailBody contains the canned body of a conductor.ListDetail response.
const ConductorListDetailBody = `
{
  "conductors": [
    {
      "links": [
        {
          "href": "http://127.0.0.1:6385/v1/conductors/compute1.localdomain",
          "rel": "self"
        },
        {
          "href": "http://127.0.0.1:6385/conductors/compute1.localdomain",
          "rel": "bookmark"
        }
      ],
      "created_at": "2018-08-07T08:39:21+00:00",
      "hostname": "compute1.localdomain",
      "conductor_group": "",
      "updated_at": "2018-11-30T07:07:23+00:00",
      "alive": false,
      "drivers": [
        "ipmi"
      ]
    },
    {
      "links": [
        {
          "href": "http://127.0.0.1:6385/v1/conductors/compute2.localdomain",
          "rel": "self"
        },
        {
          "href": "http://127.0.0.1:6385/conductors/compute2.localdomain",
          "rel": "bookmark"
        }
      ],
      "created_at": "2018-12-05T07:03:19+00:00",
      "hostname": "compute2.localdomain",
      "conductor_group": "",
      "updated_at": "2018-12-05T07:03:21+00:00",
      "alive": true,
      "drivers": [
        "ipmi"
      ]
    }
  ]
}
`

// SingleConductorBody is the canned body of a Get request on an existing conductor.
const SingleConductorBody = `
{
  "links": [
    {
      "href": "http://127.0.0.1:6385/v1/conductors/compute2.localdomain",
      "rel": "self"
    },
    {
      "href": "http://127.0.0.1:6385/conductors/compute2.localdomain",
      "rel": "bookmark"
    }
  ],
  "created_at": "2018-12-05T07:03:19+00:00",
  "hostname": "compute2.localdomain",
  "conductor_group": "",
  "updated_at": "2018-12-05T07:03:21+00:00",
  "alive": true,
  "drivers": [
    "ipmi"
  ]
}
`

var (
	createdAtFoo, _ = time.Parse(time.RFC3339, "2018-12-05T07:03:19+00:00")
	updatedAt, _    = time.Parse(time.RFC3339, "2018-12-05T07:03:21+00:00")

	ConductorFoo = conductors.Conductor{
		CreatedAt:      createdAtFoo,
		UpdatedAt:      updatedAt,
		Hostname:       "compute2.localdomain",
		ConductorGroup: "",
		Alive:          true,
		Drivers: []string{
			"ipmi",
		},
	}
)

// HandleConductorListSuccessfully sets up the test server to respond to a server List request.
func HandleConductorListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/conductors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}

		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, ConductorListBody)

		case "9e5476bd-a4ec-4653-93d6-72c93aa682ba":
			fmt.Fprintf(w, `{ "servers": [] }`)
		default:
			t.Fatalf("/conductors invoked with unexpected marker=[%s]", marker)
		}
	})
}

// HandleConductorListDetailSuccessfully sets up the test server to respond to a server List request.
func HandleConductorListDetailSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/conductors", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}

		fmt.Fprintf(w, ConductorListDetailBody)
	})
}

func HandleConductorGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/conductors/1234asdf", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, SingleConductorBody)
	})
}
