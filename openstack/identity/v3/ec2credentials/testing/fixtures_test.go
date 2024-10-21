package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2credentials"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const userID = "2844b2a08be147a08ef58317d6471f1f"
const credentialID = "f741662395b249c9b8acdebf1722c5ae"

// ListOutput provides a single page of EC2Credential results.
const ListOutput = `
{
  "credentials": [
    {
      "user_id": "2844b2a08be147a08ef58317d6471f1f",
      "links": {
        "self": "http://identity:5000/v3/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2/f741662395b249c9b8acdebf1722c5ae"
      },
      "tenant_id": "6238dee2fec940a6bf31e49e9faf995a",
      "access": "f741662395b249c9b8acdebf1722c5ae",
      "secret": "6a61eb0296034c89b49cc51dde9b40aa",
      "trust_id": null
    },
    {
      "user_id": "2844b2a08be147a08ef58317d6471f1f",
      "links": {
        "self": "http://identity:5000/v3/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2/ad6fc85fc2df49e6b5c23d5b5bdff980"
      },
      "tenant_id": "6238dee2fec940a6bf31e49e9faf995a",
      "access": "ad6fc85fc2df49e6b5c23d5b5bdff980",
      "secret": "eb233f680a204097ac329ebe8dba6d32",
      "trust_id": null
    }
  ],
  "links": {
    "self": "http://identity:5000/v3/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2",
    "previous": null,
    "next": null
  }
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
  "credential": {
    "user_id": "2844b2a08be147a08ef58317d6471f1f",
    "links": {
      "self": "http://identity:5000/v3/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2/f741662395b249c9b8acdebf1722c5ae"
    },
    "tenant_id": "6238dee2fec940a6bf31e49e9faf995a",
    "access": "f741662395b249c9b8acdebf1722c5ae",
    "secret": "6a61eb0296034c89b49cc51dde9b40aa",
    "trust_id": null
  }
}
`

// CreateRequest provides the input to a Create request.
const CreateRequest = `
{
  "tenant_id": "6238dee2fec940a6bf31e49e9faf995a"
}
`

const CreateResponse = `
{
  "credential": {
    "user_id": "2844b2a08be147a08ef58317d6471f1f",
    "links": {
      "self": "http://identity:5000/v3/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2/f741662395b249c9b8acdebf1722c5ae"
    },
    "tenant_id": "6238dee2fec940a6bf31e49e9faf995a",
    "access": "f741662395b249c9b8acdebf1722c5ae",
    "secret": "6a61eb0296034c89b49cc51dde9b40aa",
    "trust_id": null
  }
}
`

var EC2Credential = ec2credentials.Credential{
	UserID:   "2844b2a08be147a08ef58317d6471f1f",
	TenantID: "6238dee2fec940a6bf31e49e9faf995a",
	Access:   "f741662395b249c9b8acdebf1722c5ae",
	Secret:   "6a61eb0296034c89b49cc51dde9b40aa",
	TrustID:  "",
	Links: map[string]any{
		"self": "http://identity:5000/v3/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2/f741662395b249c9b8acdebf1722c5ae",
	},
}

var SecondEC2Credential = ec2credentials.Credential{
	UserID:   "2844b2a08be147a08ef58317d6471f1f",
	TenantID: "6238dee2fec940a6bf31e49e9faf995a",
	Access:   "ad6fc85fc2df49e6b5c23d5b5bdff980",
	Secret:   "eb233f680a204097ac329ebe8dba6d32",
	TrustID:  "",
	Links: map[string]any{
		"self": "http://identity:5000/v3/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2/ad6fc85fc2df49e6b5c23d5b5bdff980",
	},
}

// ExpectedEC2CredentialsSlice is the slice of application credentials expected to be returned from ListOutput.
var ExpectedEC2CredentialsSlice = []ec2credentials.Credential{EC2Credential, SecondEC2Credential}

// HandleListEC2CredentialsSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that responds with a list of two applicationcredentials.
func HandleListEC2CredentialsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetEC2CredentialSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that responds with a single application credential.
func HandleGetEC2CredentialSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2/f741662395b249c9b8acdebf1722c5ae", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleCreateEC2CredentialSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that tests application credential creation.
func HandleCreateEC2CredentialSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, CreateResponse)
	})
}

// HandleDeleteEC2CredentialSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that tests application credential deletion.
func HandleDeleteEC2CredentialSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/credentials/OS-EC2/f741662395b249c9b8acdebf1722c5ae", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}
