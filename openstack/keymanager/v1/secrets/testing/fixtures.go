package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// ListResult provides a single page of RESOURCE results.
const ListResult = `
{
    "secrets": [
        {
            "algorithm": "aes",
            "bit_length": 256,
            "content_types": {
                "default": "text/plain"
            },
            "created": "2018-06-21T02:49:48",
            "creator_id": "5c70d99f4a8641c38f8084b32b5e5c0e",
            "expiration": null,
            "mode": "cbc",
            "name": "mysecret",
            "secret_ref": "http://barbican:9311/v1/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c",
            "secret_type": "opaque",
            "status": "ACTIVE",
            "updated": "2018-06-21T02:49:48"
        },
        {
            "algorithm": "aes",
            "bit_length": 256,
            "content_types": {
                "default": "text/plain"
            },
            "created": "2018-06-21T05:18:45",
            "creator_id": "5c70d99f4a8641c38f8084b32b5e5c0e",
            "expiration": null,
            "mode": "cbc",
            "name": "anothersecret",
            "secret_ref": "http://barbican:9311/v1/secrets/1b12b69a-8822-442e-a303-da24ade648ac",
            "secret_type": "opaque",
            "status": "ACTIVE",
            "updated": "2018-06-21T05:18:45"
        }
    ],
    "total": 2
}`

// GetResult provides a Get result.
const GetResult = `
{
    "algorithm": "aes",
    "bit_length": 256,
    "content_types": {
        "default": "text/plain"
    },
    "created": "2018-06-21T02:49:48",
    "creator_id": "5c70d99f4a8641c38f8084b32b5e5c0e",
    "expiration": null,
    "mode": "cbc",
    "name": "mysecret",
    "secret_ref": "http://barbican:9311/v1/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c",
    "secret_type": "opaque",
    "status": "ACTIVE",
    "updated": "2018-06-21T02:49:48"
}`

// GetPayloadResult provides a payload result.
const GetPayloadResult = `foobar`

// CreateRequest provides the input to a Create request.
const CreateRequest = `
{
    "algorithm": "aes",
    "bit_length": 256,
    "mode": "cbc",
    "name": "mysecret",
    "payload": "foobar",
    "payload_content_type": "text/plain",
    "secret_type": "opaque"
}`

// CreateResult provides a Create result.
const CreateResult = `
{
	"secret_ref": "http://barbican:9311/v1/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c"
}`

// UpdateRequest provides the input to as Update request.
const UpdateRequest = `foobar`

// FirstSecret is the first secret in the List request.
var FirstSecret = secrets.Secret{
	Algorithm: "aes",
	BitLength: 256,
	ContentTypes: map[string]string{
		"default": "text/plain",
	},
	Created:    time.Date(2018, 6, 21, 2, 49, 48, 0, time.UTC),
	CreatorID:  "5c70d99f4a8641c38f8084b32b5e5c0e",
	Mode:       "cbc",
	Name:       "mysecret",
	SecretRef:  "http://barbican:9311/v1/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c",
	SecretType: "opaque",
	Status:     "ACTIVE",
	Updated:    time.Date(2018, 6, 21, 2, 49, 48, 0, time.UTC),
}

// SecondSecret is the second secret in the List request.
var SecondSecret = secrets.Secret{
	Algorithm: "aes",
	BitLength: 256,
	ContentTypes: map[string]string{
		"default": "text/plain",
	},
	Created:    time.Date(2018, 6, 21, 5, 18, 45, 0, time.UTC),
	CreatorID:  "5c70d99f4a8641c38f8084b32b5e5c0e",
	Mode:       "cbc",
	Name:       "anothersecret",
	SecretRef:  "http://barbican:9311/v1/secrets/1b12b69a-8822-442e-a303-da24ade648ac",
	SecretType: "opaque",
	Status:     "ACTIVE",
	Updated:    time.Date(2018, 6, 21, 5, 18, 45, 0, time.UTC),
}

// ExpectedSecretsSlice is the slice of secrets expected to be returned from ListResult.
var ExpectedSecretsSlice = []secrets.Secret{FirstSecret, SecondSecret}

// ExpectedCreateResult is the result of a create request
var ExpectedCreateResult = secrets.Secret{
	SecretRef: "http://barbican:9311/v1/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c",
}

// HandleListSecretsSuccessfully creates an HTTP handler at `/secrets` on the
// test handler mux that responds with a list of two secrets.
func HandleListSecretsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/secrets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListResult)
	})
}

// HandleGetSecretSuccessfully creates an HTTP handler at `/secrets` on the
// test handler mux that responds with a single secret.
func HandleGetSecretSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetResult)
	})
}

// HandleGetPayloadSuccessfully creates an HTTP handler at `/secrets` on the
// test handler mux that responds with a single secret.
func HandleGetPayloadSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c/payload", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetPayloadResult)
	})
}

// HandleCreateSecretSuccessfully creates an HTTP handler at `/secrets` on the
// test handler mux that tests secret creation.
func HandleCreateSecretSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/secrets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreateResult)
	})
}

// HandleDeleteSecretSuccessfully creates an HTTP handler at `/secrets` on the
// test handler mux that tests secret deletion.
func HandleDeleteSecretSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleUpdateSecretSuccessfully creates an HTTP handler at `/secrets` on the
// test handler mux that tests secret updates.
func HandleUpdateSecretSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/secrets/1b8068c4-3bb6-4be6-8f1e-da0d1ea0b67c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNoContent)
	})
}
