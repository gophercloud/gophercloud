package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/tsigkeys"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

// ListOutput is a sample response to a List call.
const ListOutput = `
{
    "links": {
        "self": "http://example.com:9001/v2/tsigkeys"
    },
    "metadata": {
        "total_count": 2
    },
    "tsigkeys": [
        {
            "id": "8add45a3-0f29-489f-854e-7609baf8d7a1",
            "name": "poolsecondarykey",
            "algorithm": "hmac-sha256",
            "secret": "my-base64-secret-example==",
            "scope": "POOL",
            "resource_id": "adcc2fb6-7984-4453-a6f9-2cc2a24a38bb",
            "created_at": "2025-08-13T15:54:18.000000",
            "updated_at": null,
            "links": {
                "self": "http://127.0.0.1:9001/v2/tsigkeys/8add45a3-0f29-489f-854e-7609baf8d7a1"
            }
        },
        {
            "id": "9bef46b4-1f3a-59a0-965f-8710caf9e8b2",
            "name": "zonekey",
            "algorithm": "hmac-sha512",
            "secret": "another-base64-secret-example==",
            "scope": "ZONE",
            "resource_id": "c4d5e6f7-8a9b-0c1d-2e3f-4a5b6c7d8e9f",
            "created_at": "2025-08-14T10:30:45.000000",
            "updated_at": "2025-08-15T14:22:33.000000",
            "links": {
                "self": "http://127.0.0.1:9001/v2/tsigkeys/9bef46b4-1f3a-59a0-965f-8710caf9e8b2"
            }
        }
    ]
}
`

// GetOutput is a sample response to a Get call.
const GetOutput = `
{
    "id": "8add45a3-0f29-489f-854e-7609baf8d7a1",
    "name": "poolsecondarykey",
    "algorithm": "hmac-sha256",
    "secret": "my-base64-secret-example==",
    "scope": "POOL",
    "resource_id": "adcc2fb6-7984-4453-a6f9-2cc2a24a38bb",
    "created_at": "2025-08-13T15:54:18.000000",
    "updated_at": null,
    "links": {
        "self": "http://127.0.0.1:9001/v2/tsigkeys/8add45a3-0f29-489f-854e-7609baf8d7a1"
    }
}
`

// FirstTSIGKeyCreatedAt is the created at time for the first TSIG key
var FirstTSIGKeyCreatedAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2025-08-13T15:54:18.000000")

// FirstTSIGKey is the first result in ListOutput
var FirstTSIGKey = tsigkeys.TSIGKey{
	ID:         "8add45a3-0f29-489f-854e-7609baf8d7a1",
	Name:       "poolsecondarykey",
	Algorithm:  "hmac-sha256",
	Secret:     "my-base64-secret-example==",
	Scope:      "POOL",
	ResourceID: "adcc2fb6-7984-4453-a6f9-2cc2a24a38bb",
	CreatedAt:  FirstTSIGKeyCreatedAt,
	Links: map[string]any{
		"self": "http://127.0.0.1:9001/v2/tsigkeys/8add45a3-0f29-489f-854e-7609baf8d7a1",
	},
}

var SecondTSIGKeyCreatedAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2025-08-14T10:30:45.000000")
var SecondTSIGKeyUpdatedAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2025-08-15T14:22:33.000000")

var SecondTSIGKey = tsigkeys.TSIGKey{
	ID:         "9bef46b4-1f3a-59a0-965f-8710caf9e8b2",
	Name:       "zonekey",
	Algorithm:  "hmac-sha512",
	Secret:     "another-base64-secret-example==",
	Scope:      "ZONE",
	ResourceID: "c4d5e6f7-8a9b-0c1d-2e3f-4a5b6c7d8e9f",
	CreatedAt:  SecondTSIGKeyCreatedAt,
	UpdatedAt:  SecondTSIGKeyUpdatedAt,
	Links: map[string]any{
		"self": "http://127.0.0.1:9001/v2/tsigkeys/9bef46b4-1f3a-59a0-965f-8710caf9e8b2",
	},
}

// ExpectedTSIGKeysSlice is the slice of results that should be parsed
// from ListOutput, in the expected order.
var ExpectedTSIGKeysSlice = []tsigkeys.TSIGKey{FirstTSIGKey, SecondTSIGKey}

// HandleListSuccessfully configures the test server to respond to a List request.
func HandleListSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/tsigkeys", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, ListOutput)
	})
}

// HandleGetSuccessfully configures the test server to respond to a Get request.
func HandleGetSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/tsigkeys/8add45a3-0f29-489f-854e-7609baf8d7a1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, GetOutput)
	})
}

// CreateTSIGKeyRequest is a sample request to create a TSIG key.
const CreateTSIGKeyRequest = `
{
    "name": "poolsecondarykey",
    "algorithm": "hmac-sha256",
    "secret": "my-base64-secret-example==",
    "scope": "POOL",
    "resource_id": "adcc2fb6-7984-4453-a6f9-2cc2a24a38bb"
}
`

// CreateTSIGKeyResponse is a sample response to a create request.
const CreateTSIGKeyResponse = `
{
    "id": "8add45a3-0f29-489f-854e-7609baf8d7a1",
    "name": "poolsecondarykey",
    "algorithm": "hmac-sha256",
    "secret": "my-base64-secret-example==",
    "scope": "POOL",
    "resource_id": "adcc2fb6-7984-4453-a6f9-2cc2a24a38bb",
    "created_at": "2025-08-13T15:54:18.000000",
    "updated_at": null,
    "links": {
        "self": "http://127.0.0.1:9001/v2/tsigkeys/8add45a3-0f29-489f-854e-7609baf8d7a1"
    }
}
`

// CreatedTSIGKey is the expected created TSIG key
var CreatedTSIGKey = FirstTSIGKey

// HandleCreateSuccessfully configures the test server to respond to a Create request.
func HandleCreateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/tsigkeys", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateTSIGKeyRequest)

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, CreateTSIGKeyResponse)
	})
}

// UpdateTSIGKeyRequest is a sample request to update a TSIG key.
const UpdateTSIGKeyRequest = `
{
    "name": "updatedsecondarykey",
    "secret": "new-base64-secret-example=="
}
`

// UpdateTSIGKeyResponse is a sample response to update a TSIG key.
const UpdateTSIGKeyResponse = `
{
    "id": "8add45a3-0f29-489f-854e-7609baf8d7a1",
    "name": "updatedsecondarykey",
    "algorithm": "hmac-sha256",
    "secret": "new-base64-secret-example==",
    "scope": "POOL",
    "resource_id": "adcc2fb6-7984-4453-a6f9-2cc2a24a38bb",
    "created_at": "2025-08-13T15:54:18.000000",
    "updated_at": "2025-08-16T09:15:22.000000",
    "links": {
        "self": "http://127.0.0.1:9001/v2/tsigkeys/8add45a3-0f29-489f-854e-7609baf8d7a1"
    }
}
`

// HandleUpdateSuccessfully configures the test server to respond to an Update request.
func HandleUpdateSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/tsigkeys/8add45a3-0f29-489f-854e-7609baf8d7a1",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PATCH")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, UpdateTSIGKeyRequest)

			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-Type", "application/json")
			fmt.Fprint(w, UpdateTSIGKeyResponse)
		})
}

// HandleDeleteSuccessfully configures the test server to respond to a Delete request.
func HandleDeleteSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/tsigkeys/8add45a3-0f29-489f-854e-7609baf8d7a1",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "DELETE")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

			w.WriteHeader(http.StatusNoContent)
		})
}
