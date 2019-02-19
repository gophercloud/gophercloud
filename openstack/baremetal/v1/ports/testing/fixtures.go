package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// PortListBody contains the canned body of a ports.List response, without detail.
const PortListBody = `
{
   "ports": [
       {
           "uuid": "3abe3f36-9708-4e9f-b07e-0f898061d3a7",
           "links": [
               {
                   "href": "http://192.168.0.8/baremetal/v1/ports/3abe3f36-9708-4e9f-b07e-0f898061d3a7",
                   "rel": "self"
               },
               {
                   "href": "http://192.168.0.8/baremetal/ports/3abe3f36-9708-4e9f-b07e-0f898061d3a7",
                   "rel": "bookmark"
               }
           ],
           "address": "52:54:00:0a:af:d1"
       },
       {
           "uuid": "f2845e11-dbd4-4728-a8c0-30d19f48924a",
           "links": [
               {
                   "href": "http://192.168.0.8/baremetal/v1/ports/f2845e11-dbd4-4728-a8c0-30d19f48924a",
                   "rel": "self"
               },
               {
                   "href": "http://192.168.0.8/baremetal/ports/f2845e11-dbd4-4728-a8c0-30d19f48924a",
                   "rel": "bookmark"
               }
           ],
           "address": "52:54:00:4d:87:e6"
       }
   ]
}
`

// PortListDetailBody contains the canned body of a port.ListDetail response.
const PortListDetailBody = `
{
   "ports": [
       {
           "local_link_connection": {

           },
           "node_uuid": "ddd06a60-b91e-4ab4-a6e7-56c0b25b6086",
           "uuid": "3abe3f36-9708-4e9f-b07e-0f898061d3a7",
           "links": [
               {
                   "href": "http://192.168.0.8/baremetal/v1/ports/3abe3f36-9708-4e9f-b07e-0f898061d3a7",
                   "rel": "self"
               },
               {
                   "href": "http://192.168.0.8/baremetal/ports/3abe3f36-9708-4e9f-b07e-0f898061d3a7",
                   "rel": "bookmark"
               }
           ],
           "extra": {

           },
           "pxe_enabled": true,
           "portgroup_uuid": null,
           "updated_at": "2019-02-15T09:55:19+00:00",
           "physical_network": null,
           "address": "52:54:00:0a:af:d1",
           "internal_info": {

           },
           "created_at": "2019-02-15T09:52:23+00:00"
       },
       {
           "local_link_connection": {

           },
           "node_uuid": "ddd06a60-b91e-4ab4-a6e7-56c0b25b6086",
           "uuid": "f2845e11-dbd4-4728-a8c0-30d19f48924a",
           "links": [
               {
                   "href": "http://192.168.0.8/baremetal/v1/ports/f2845e11-dbd4-4728-a8c0-30d19f48924a",
                   "rel": "self"
               },
               {
                   "href": "http://192.168.0.8/baremetal/ports/f2845e11-dbd4-4728-a8c0-30d19f48924a",
                   "rel": "bookmark"
               }
           ],
           "extra": {

           },
           "pxe_enabled": true,
           "portgroup_uuid": null,
           "updated_at": "2019-02-15T09:55:19+00:00",
           "physical_network": null,
           "address": "52:54:00:4d:87:e6",
           "internal_info": {

           },
           "created_at": "2019-02-15T09:52:24+00:00"
       }
   ]
}
`

// SinglePortBody is the canned body of a Get request on an existing port.
const SinglePortBody = `
{
    "local_link_connection": {

    },
    "node_uuid": "ddd06a60-b91e-4ab4-a6e7-56c0b25b6086",
    "uuid": "f2845e11-dbd4-4728-a8c0-30d19f48924a",
    "links": [
        {
            "href": "http://192.168.0.8/baremetal/v1/ports/f2845e11-dbd4-4728-a8c0-30d19f48924a",
            "rel": "self"
        },
        {
            "href": "http://192.168.0.8/baremetal/ports/f2845e11-dbd4-4728-a8c0-30d19f48924a",
            "rel": "bookmark"
        }
    ],
    "extra": {

    },
    "pxe_enabled": true,
    "portgroup_uuid": null,
    "updated_at": "2019-02-15T09:55:19+00:00",
    "physical_network": null,
    "address": "52:54:00:4d:87:e6",
    "internal_info": {

    },
    "created_at": "2019-02-15T09:52:24+00:00"
}
`

var (
	PortFoo = ports.Port{
		UUID:       "f2845e11-dbd4-4728-a8c0-30d19f48924a",
		NodeUUID:   "ddd06a60-b91e-4ab4-a6e7-56c0b25b6086",
		Address:    "52:54:00:4d:87:e6",
		PXEEnabled: true,
	}

	PortBar = ports.Port{
		UUID:       "3abe3f36-9708-4e9f-b07e-0f898061d3a7",
		NodeUUID:   "ddd06a60-b91e-4ab4-a6e7-56c0b25b6086",
		Address:    "52:54:00:0a:af:d1",
		PXEEnabled: true,
	}
)

// HandlePortListSuccessfully sets up the test server to respond to a port List request.
func HandlePortListSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()

		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, PortListBody)

		case "f2845e11-dbd4-4728-a8c0-30d19f48924a":
			fmt.Fprintf(w, `{ "ports": [] }`)
		default:
			t.Fatalf("/ports invoked with unexpected marker=[%s]", marker)
		}
	})
}

// HandlePortListSuccessfully sets up the test server to respond to a server List request.
func HandlePortListDetailSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/ports/detail", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		r.ParseForm()

		fmt.Fprintf(w, PortListDetailBody)
	})
}

// HandleSPortCreationSuccessfully sets up the test server to respond to a server creation request
// with a given response.
func HandlePortCreationSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/v1/ports", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
          "boot_interface": "pxe",
          "driver": "ipmi",
          "name": "foo"
        }`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, response)
	})
}

// HandlePortDeletionSuccessfully sets up the test server to respond to a port deletion request.
func HandlePortDeletionSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/ports/3abe3f36-9708-4e9f-b07e-0f898061d3a7", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandlePortGetSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/v1/ports/f2845e11-dbd4-4728-a8c0-30d19f48924a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")

		fmt.Fprintf(w, SinglePortBody)
	})
}

func HandlePortUpdateSuccessfully(t *testing.T, response string) {
	th.Mux.HandleFunc("/v1/ports/f2845e11-dbd4-4728-a8c0-30d19f48924a", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestJSONRequest(t, r, `[{"op": "replace", "path": "/address", "value": "22:22:22:22:22:22"}]`)

		fmt.Fprintf(w, response)
	})
}
