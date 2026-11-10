package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/portgroups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// PortGroupsListBody is the JSON response for listing all portgroups.
var PortGroupsListBody = `
{
  "portgroups": [
    {
      "address": "00:1a:2b:3c:4d:5e",
      "created_at": "2019-02-20T09:43:58Z",
      "extra": {
        "description": "Primary network bond",
        "location": "rack-3-unit-12"
      },
      "internal_info": {
        "fault_count": 0,
        "last_check": "2024-03-15T10:30:00Z"
      },
      "links": [
        {
          "href": "http://ironic.example.com/v1/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796",
          "rel": "bookmark"
        }
      ],
      "mode": "active-backup",
      "name": "bond0",
      "node_uuid": "f9c9a846-c53f-4b17-9f0c-dd9f459d35c8",
      "ports": [
        {
          "href": "http://ironic.example.com/v1/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796/ports",
          "rel": "self"
        }
      ],
      "properties": {
        "miimon": "100",
        "updelay": "1000",
        "downdelay": "1000",
        "xmit_hash_policy": "layer2"
      },
      "standalone_ports_supported": true,
      "updated_at": "2019-02-20T09:43:58Z",
      "uuid": "d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796"
    },
    {
      "address": "11:22:33:44:55:66",
      "created_at": "2019-02-20T09:43:58Z",
      "extra": {
        "description": "Secondary bond",
        "location": "rack-1-unit-4"
      },
      "internal_info": {
        "fault_count": 1,
        "last_check": "2024-04-01T09:00:00Z"
      },
      "links": [
        {
          "href": "http://ironic.example.com/v1/portgroups/a1b2c3d4-e5f6-7890-1234-56789abcdef0",
          "rel": "self"
        },
        {
          "href": "http://ironic.example.com/portgroups/a1b2c3d4-e5f6-7890-1234-56789abcdef0",
          "rel": "bookmark"
        }
      ],
      "mode": "active-backup",
      "name": "bond1",
      "node_uuid": "aabbcc00-1122-3344-5566-778899aabbcc",
      "ports": [
        {
          "href": "http://ironic.example.com/v1/portgroups/a1b2c3d4-e5f6-7890-1234-56789abcdef0/ports",
          "rel": "self"
        }
      ],
      "properties": {
        "miimon": "200",
        "updelay": "500",
        "downdelay": "500",
        "xmit_hash_policy": "layer3+4"
      },
      "standalone_ports_supported": true,
      "updated_at": "2019-02-20T09:43:58Z",
      "uuid": "a1b2c3d4-e5f6-7890-1234-56789abcdef0"
    }
  ]
}
`

// SinglePortGroupBody returns JSON for a single portgroup.
// Here we use PortGroup1 as the example.
var SinglePortGroupBody = `
{
  "address": "00:1a:2b:3c:4d:5e",
  "created_at": "2019-02-20T09:43:58Z",
  "extra": {
    "description": "Primary network bond",
    "location": "rack-3-unit-12"
  },
  "internal_info": {
    "fault_count": 0,
    "last_check": "2024-03-15T10:30:00Z"
  },
  "links": [
    {
      "href": "http://ironic.example.com/v1/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796",
      "rel": "self"
    },
    {
      "href": "http://ironic.example.com/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796",
      "rel": "bookmark"
    }
  ],
  "mode": "active-backup",
  "name": "bond0",
  "node_uuid": "f9c9a846-c53f-4b17-9f0c-dd9f459d35c8",
  "ports": [
    {
      "href": "http://ironic.example.com/v1/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796/ports",
      "rel": "self"
    }
  ],
  "properties": {
    "miimon": "100",
    "updelay": "1000",
    "downdelay": "1000",
    "xmit_hash_policy": "layer2"
  },
  "standalone_ports_supported": true,
  "updated_at": "2019-02-20T09:43:58Z",
  "uuid": "d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796"
}
`

var (
	createdAt, _ = time.Parse(time.RFC3339, "2019-02-20T09:43:58Z")

	// PortGroup1 is the first portgroup.
	PortGroup1 = portgroups.PortGroup{
		Name:                     "bond0",
		UUID:                     "d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796",
		Address:                  "00:1a:2b:3c:4d:5e",
		NodeUUID:                 "f9c9a846-c53f-4b17-9f0c-dd9f459d35c8",
		StandalonePortsSupported: true,
		InternalInfo: map[string]any{
			"fault_count": float64(0),
			"last_check":  "2024-03-15T10:30:00Z",
		},
		Extra: map[string]any{
			"description": "Primary network bond",
			"location":    "rack-3-unit-12",
		},
		Mode: "active-backup",
		Properties: map[string]any{
			"miimon":           "100",
			"updelay":          "1000",
			"downdelay":        "1000",
			"xmit_hash_policy": "layer2",
		},
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		Links: []portgroups.ResourceLink{
			{
				Href: "http://ironic.example.com/v1/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796",
				Rel:  "self",
			},
			{
				Href: "http://ironic.example.com/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796",
				Rel:  "bookmark",
			},
		},
		Ports: []portgroups.ResourceLink{
			{
				Href: "http://ironic.example.com/v1/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796/ports",
				Rel:  "self",
			},
		},
	}
)

// HandlePortGroupListSuccessfully sets up the test server to respond to a
// portgroup List request.
func HandlePortGroupListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/portgroups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form: %v", err)
		}

		marker := r.Form.Get("marker")
		switch marker {
		case "":
			// Return both portgroups.
			fmt.Fprint(w, PortGroupsListBody)
		case "d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796":
			// No portgroups remain.
			fmt.Fprintf(w, `{ "portgroups": [] }`)
		default:
			t.Fatalf("/portgroups invoked with unexpected marker=[%s]", marker)
		}
	})
}

// HandlePortGroupCreationSuccessfully sets up the test server to respond to a PortGroup creation request
// with a given response.
func HandlePortGroupCreationSuccessfully(t *testing.T, fakeServer th.FakeServer, response string) {
	fakeServer.Mux.HandleFunc("/portgroups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `{
      "node_uuid": "f9c9a846-c53f-4b17-9f0c-dd9f459d35c8",
      "address": "00:1a:2b:3c:4d:5e",
      "name": "bond0",
      "mode": "active-backup",
      "standalone_ports_supported": true,
      "properties": {
          "miimon": "100",
          "updelay": "1000",
          "downdelay": "1000",
          "xmit_hash_policy": "layer2"
      },
      "extra": {
          "description": "Primary network bond",
          "location": "rack-3-unit-12"
      }
  }`)

		w.WriteHeader(http.StatusAccepted)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, response)
	})
}

// HandlePortGroupDeletionSuccessfully sets up the test server to respond to a
// portgroup Deletion (DELETE) request for PortGroup2.
func HandlePortGroupDeletionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			w.WriteHeader(http.StatusNoContent)
		})
}

// HandlePortGroupGetSuccessfully sets up the test server to respond to a
// portgroup Get request for PortGroup1.
func HandlePortGroupGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/portgroups/d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "GET")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestHeader(t, r, "Accept", "application/json")

			w.Header().Add("Content-Type", "application/json")
			fmt.Fprint(w, SinglePortGroupBody)
		})
}
