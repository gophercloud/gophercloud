package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListOutput is a sample response to a List call.
const ListOutput = `
{
    "server_groups": [
        {
            "id": "616fb98f-46ca-475e-917e-2563e5a8cd19",
            "name": "test",
            "policies": [
                "anti-affinity"
            ],
            "members": [],
            "metadata": {}
        },
        {
            "id": "4d8c3732-a248-40ed-bebc-539a6ffd25c0",
            "name": "test2",
            "policies": [
                "affinity"
            ],
            "members": [],
            "metadata": {}
        }
    ]
}
`

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
    "server_group": {
        "id": "616fb98f-46ca-475e-917e-2563e5a8cd19",
        "name": "test",
        "policies": [
            "anti-affinity"
        ],
        "members": [],
        "metadata": {}
    }
}
`

// GetOutputMicroversion is a sample response to a Get call with microversion set to 2.64
const GetOutputMicroversion = `
{
    "server_group": {
        "id": "616fb98f-46ca-475e-917e-2563e5a8cd19",
        "name": "test",
        "policy": "anti-affinity",
        "rules": {
          "max_server_per_host": 3
        },
        "members": [],
        "metadata": {}
    }
}
`

// CreateOutput is a sample response to a Post call
const CreateOutput = `
{
    "server_group": {
        "id": "616fb98f-46ca-475e-917e-2563e5a8cd19",
        "name": "test",
        "policies": [
            "anti-affinity"
        ],
        "members": [],
        "metadata": {}
    }
}
`

// CreateOutputMicroversion is a sample response to a Post call with microversion set to 2.64
const CreateOutputMicroversion = `
{
    "server_group": {
        "id": "616fb98f-46ca-475e-917e-2563e5a8cd19",
        "name": "test",
        "policy": "anti-affinity",
        "rules": {
          "max_server_per_host": 3
        },
        "members": [],
        "metadata": {}
    }
}
`

var policy = "anti-affinity"

// ExpectedServerGroupGet is parsed result from GetOutput.
var ExpectedServerGroupGet = servergroups.ServerGroup{
	ID:   "616fb98f-46ca-475e-917e-2563e5a8cd19",
	Name: "test",
	Policies: []string{
		"anti-affinity",
	},
	Members:  []string{},
	Metadata: map[string]any{},
}

// ExpectedServerGroupGet is parsed result from GetOutputMicroversion.
var ExpectedServerGroupGetMicroversion = servergroups.ServerGroup{
	ID:     "616fb98f-46ca-475e-917e-2563e5a8cd19",
	Name:   "test",
	Policy: &policy,
	Rules: &servergroups.Rules{
		MaxServerPerHost: 3,
	},
	Members:  []string{},
	Metadata: map[string]any{},
}

// ExpectedServerGroupList is the slice of results that should be parsed
// from ListOutput, in the expected order.
var ExpectedServerGroupList = []servergroups.ServerGroup{
	{
		ID:   "616fb98f-46ca-475e-917e-2563e5a8cd19",
		Name: "test",
		Policies: []string{
			"anti-affinity",
		},
		Members:  []string{},
		Metadata: map[string]any{},
	},
	{
		ID:   "4d8c3732-a248-40ed-bebc-539a6ffd25c0",
		Name: "test2",
		Policies: []string{
			"affinity",
		},
		Members:  []string{},
		Metadata: map[string]any{},
	},
}

// ExpectedServerGroupCreate is the parsed result from CreateOutput.
var ExpectedServerGroupCreate = servergroups.ServerGroup{
	ID:   "616fb98f-46ca-475e-917e-2563e5a8cd19",
	Name: "test",
	Policies: []string{
		"anti-affinity",
	},
	Members:  []string{},
	Metadata: map[string]any{},
}

// CreatedServerGroup is the parsed result from CreateOutputMicroversion.
var ExpectedServerGroupCreateMicroversion = servergroups.ServerGroup{
	ID:     "616fb98f-46ca-475e-917e-2563e5a8cd19",
	Name:   "test",
	Policy: &policy,
	Rules: &servergroups.Rules{
		MaxServerPerHost: 3,
	},
	Members:  []string{},
	Metadata: map[string]any{},
}

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-server-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, ListOutput)
	})
}

// HandleGetSuccessfully configures the test server to respond to a Get request
// for an existing server group
func HandleGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-server-groups/4d8c3732-a248-40ed-bebc-539a6ffd25c0", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, GetOutput)
	})
}

// HandleGetMicroversionSuccessfully configures the test server to respond to a Get request
// for an existing server group with microversion set to 2.64
func HandleGetMicroversionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-server-groups/4d8c3732-a248-40ed-bebc-539a6ffd25c0", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, GetOutputMicroversion)
	})
}

// HandleCreateSuccessfully configures the test server to respond to a Create request
// for a new server group
func HandleCreateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-server-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
{
    "server_group": {
        "name": "test",
        "policies": [
            "anti-affinity"
        ]
    }
}
`)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, CreateOutput)
	})
}

// HandleCreateMicroversionSuccessfully configures the test server to respond to a Create request
// for a new server group with microversion set to 2.64
func HandleCreateMicroversionSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-server-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
{
    "server_group": {
        "name": "test",
        "policy": "anti-affinity",
        "rules": {
            "max_server_per_host": 3
        }
    }
}
`)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, CreateOutputMicroversion)
	})
}

// HandleDeleteSuccessfully configures the test server to respond to a Delete request for a
// an existing server group
func HandleDeleteSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/os-server-groups/616fb98f-46ca-475e-917e-2563e5a8cd19", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})
}
