package testing

import (
	"fmt"
	"net/http"
	"testing"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListOutput is a sample response to a List call.
const ListOutput = `
{
  "shares": [
    {
      "share_id": "e8debdc0-447a-4376-a10a-4cd9122d7986",
      "status": "inactive",
      "tag": "e8debdc0-447a-4376-a10a-4cd9122d7986"
    },
    {
      "share_id": "a26887c6-c47b-4654-abb5-dfadf7d3f803",
      "status": "active",
      "tag": "a26887c6-c47b-4654-abb5-dfadf7d3f803"
    }
  ]
}
`

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
  "share": {
    "share_id": "e8debdc0-447a-4376-a10a-4cd9122d7986",
    "status": "inactive",
    "tag": "e8debdc0-447a-4376-a10a-4cd9122d7986"
  }
}
`

// GetOutputAdmin is a sample response to a Get call for an admin.
const GetOutputAdmin = `
{
  "share": {
    "uuid": "68ba1762-fd6d-4221-8311-f3193dd93404",
    "share_id": "e8debdc0-447a-4376-a10a-4cd9122d7986",
    "status": "inactive",
    "export_location": "10.0.0.50:/mnt/foo",
    "tag": "e8debdc0-447a-4376-a10a-4cd9122d7986"
  }
}
`

// CreateOutput is a sample response to a Create call.
const CreateOutput = `
{
  "share": {
    "share_id": "e8debdc0-447a-4376-a10a-4cd9122d7986",
    "status": "attaching",
    "tag": "my-share"
  }
}
`

// CreateOutputWithoutTag is a sample response to a Create call when the request
// omits tag. Nova returns tag equal to share_id.
const CreateOutputWithoutTag = `
{
  "share": {
    "share_id": "e8debdc0-447a-4376-a10a-4cd9122d7986",
    "status": "attaching",
    "tag": "e8debdc0-447a-4376-a10a-4cd9122d7986"
  }
}
`

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/servers/4d8c3732-a248-40ed-bebc-539a6ffd25c0/shares", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, ListOutput)
	})
}

// HandleGetSuccessfully configures the test server to respond to a Get request
// for an existing attachment
func HandleGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/servers/4d8c3732-a248-40ed-bebc-539a6ffd25c0/shares/e8debdc0-447a-4376-a10a-4cd9122d7986", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, GetOutput)
	})
}

// HandleGetSuccessfullyAdmin configures the test server to respond to a Get request
// for an existing attachment for an admin
func HandleGetSuccessfullyAdmin(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/servers/4d8c3732-a248-40ed-bebc-539a6ffd25c0/shares/e8debdc0-447a-4376-a10a-4cd9122d7986", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, GetOutputAdmin)
	})
}

// HandleCreateSuccessfully configures the test server to respond to a Create request
// for a new attachment
func HandleCreateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/servers/4d8c3732-a248-40ed-bebc-539a6ffd25c0/shares", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
{
    "share": {
        "share_id": "3cdf5132-64f2-4584-876a-bd296ae7eabd",
        "tag": "my-share"
    }
}
`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, CreateOutput)
	})
}

// HandleCreateSuccessfullyWithoutTag configures the test server to respond to a Create
// request that omits tag in the request body.
func HandleCreateSuccessfullyWithoutTag(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/servers/4d8c3732-a248-40ed-bebc-539a6ffd25c0/shares", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
{
    "share": {
        "share_id": "3cdf5132-64f2-4584-876a-bd296ae7eabd"
    }
}
`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, CreateOutputWithoutTag)
	})
}

// HandleDeleteSuccessfully configures the test server to respond to a Delete request for
// an existing attachment
func HandleDeleteSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/servers/4d8c3732-a248-40ed-bebc-539a6ffd25c0/shares/e8debdc0-447a-4376-a10a-4cd9122d7986", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusOK)
	})
}
