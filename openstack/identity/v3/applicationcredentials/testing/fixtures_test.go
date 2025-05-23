package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/applicationcredentials"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

const userID = "2844b2a08be147a08ef58317d6471f1f"
const applicationCredentialID = "f741662395b249c9b8acdebf1722c5ae"
const accessRuleID = "07d719df00f349ef8de77d542edf010c"

// ListOutput provides a single page of ApplicationCredential results.
const ListOutput = `
{
  "links": {
    "self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials",
    "previous": null,
    "next": null
  },
  "application_credentials": [
    {
      "links": {
        "self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/c4859fb437df4b87a51a8f5adcfb0bc7"
      },
      "description": null,
      "roles": [
        {
          "id": "31f87923ae4a4d119aa0b85dcdbeed13",
          "domain_id": null,
          "name": "compute_viewer"
        }
      ],
      "expires_at": null,
      "unrestricted": false,
      "project_id": "53c2b94f63fb4f43a21b92d119ce549f",
      "id": "c4859fb437df4b87a51a8f5adcfb0bc7",
      "name": "test1"
    },
    {
      "links": {
        "self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/6b8cc7647da64166a4a3cc0c88ebbabb"
      },
      "description": null,
      "roles": [
        {
          "id": "31f87923ae4a4d119aa0b85dcdbeed13",
          "domain_id": null,
          "name": "compute_viewer"
        },
        {
          "id": "4494bc5bea1a4105ad7fbba6a7eb9ef4",
          "domain_id": null,
          "name": "network_viewer"
        }
      ],
      "expires_at": "2019-03-12T12:12:12.123456",
      "unrestricted": true,
      "project_id": "53c2b94f63fb4f43a21b92d119ce549f",
      "id": "6b8cc7647da64166a4a3cc0c88ebbabb",
      "name": "test2"
    }
  ]
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
  "application_credential": {
    "links": {
      "self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/f741662395b249c9b8acdebf1722c5ae"
    },
    "description": null,
    "roles": [
      {
        "id": "31f87923ae4a4d119aa0b85dcdbeed13",
        "domain_id": null,
        "name": "compute_viewer"
      }
    ],
    "access_rules": [
      {
        "path": "/v2.0/metrics",
        "id": "07d719df00f349ef8de77d542edf010c",
        "service": "monitoring",
        "method": "GET"
      }
    ],
    "expires_at": null,
    "unrestricted": false,
    "project_id": "53c2b94f63fb4f43a21b92d119ce549f",
    "id": "f741662395b249c9b8acdebf1722c5ae",
    "name": "test"
  }
}
`

// CreateRequest provides the input to a Create request.
const CreateRequest = `
{
  "application_credential": {
    "secret": "mysecret",
    "unrestricted": false,
    "roles": [
      {
        "id": "31f87923ae4a4d119aa0b85dcdbeed13"
      }
    ],
    "access_rules": [
      {
        "path": "/v2.0/metrics",
        "method": "GET",
        "service": "monitoring"
      }
    ],
    "name": "test"
  }
}
`

const CreateResponse = `
{
  "application_credential": {
    "links": {
      "self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/f741662395b249c9b8acdebf1722c5ae"
    },
    "description": null,
    "roles": [
      {
        "id": "31f87923ae4a4d119aa0b85dcdbeed13",
        "domain_id": null,
        "name": "compute_viewer"
      }
    ],
    "access_rules": [
      {
        "path": "/v2.0/metrics",
        "id": "07d719df00f349ef8de77d542edf010c",
        "service": "monitoring",
        "method": "GET"
      }
    ],
    "expires_at": null,
    "secret": "mysecret",
    "unrestricted": false,
    "project_id": "53c2b94f63fb4f43a21b92d119ce549f",
    "id": "f741662395b249c9b8acdebf1722c5ae",
    "name": "test"
  }
}
`

// CreateNoOptionsRequest provides the input to a Create request with no Secret.
const CreateNoSecretRequest = `
{
  "application_credential": {
    "unrestricted": false,
    "name": "test1",
    "roles": [
      {
        "id": "31f87923ae4a4d119aa0b85dcdbeed13"
      }
    ]
  }
}
`

const CreateNoSecretResponse = `
{
  "application_credential": {
    "links": {
      "self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/c4859fb437df4b87a51a8f5adcfb0bc7"
    },
    "description": null,
    "roles": [
      {
        "id": "31f87923ae4a4d119aa0b85dcdbeed13",
        "domain_id": null,
        "name": "compute_viewer"
      }
    ],
    "expires_at": null,
    "secret": "generated_secret",
    "unrestricted": false,
    "project_id": "53c2b94f63fb4f43a21b92d119ce549f",
    "id": "c4859fb437df4b87a51a8f5adcfb0bc7",
    "name": "test1"
  }
}
`

const CreateUnrestrictedRequest = `
{
  "application_credential": {
    "unrestricted": true,
    "roles": [
      {
        "id": "31f87923ae4a4d119aa0b85dcdbeed13"
      },
      {
        "id": "4494bc5bea1a4105ad7fbba6a7eb9ef4"
      }
    ],
    "expires_at": "2019-03-12T12:12:12.123456",
    "name": "test2"
  }
}
`

const CreateUnrestrictedResponse = `
{
  "application_credential": {
    "links": {
      "self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/6b8cc7647da64166a4a3cc0c88ebbabb"
    },
    "description": null,
    "roles": [
      {
        "id": "31f87923ae4a4d119aa0b85dcdbeed13",
        "domain_id": null,
        "name": "compute_viewer"
      },
      {
        "id": "4494bc5bea1a4105ad7fbba6a7eb9ef4",
        "domain_id": null,
        "name": "network_viewer"
      }
    ],
    "expires_at": "2019-03-12T12:12:12.123456",
    "secret": "generated_secret",
    "unrestricted": true,
    "project_id": "53c2b94f63fb4f43a21b92d119ce549f",
    "id": "6b8cc7647da64166a4a3cc0c88ebbabb",
    "name": "test2"
  }
}
`

// ListAccessRulesOutput provides a single page of AccessRules results.
const ListAccessRulesOutput = `
{
  "links": {
    "self": "https://example.com/identity/v3/users/2844b2a08be147a08ef58317d6471f1f/access_rules",
    "previous": null,
    "next": null
  },
  "access_rules": [
    {
      "path": "/v2.0/metrics",
      "links": {
        "self": "https://example.com/identity/v3/access_rules/07d719df00f349ef8de77d542edf010c"
      },
      "id": "07d719df00f349ef8de77d542edf010c",
      "service": "monitoring",
      "method": "GET"
    }
  ]
}`

// GetAccessRuleOutput provides a Get result.
const GetAccessRuleOutput = `
{
  "access_rule": {
    "path": "/v2.0/metrics",
    "links": {
      "self": "https://example.com/identity/v3/access_rules/07d719df00f349ef8de77d542edf010c"
    },
    "id": "07d719df00f349ef8de77d542edf010c",
    "service": "monitoring",
    "method": "GET"
  }
}
`

var nilTime time.Time
var ApplicationCredential = applicationcredentials.ApplicationCredential{
	ID:           "f741662395b249c9b8acdebf1722c5ae",
	Name:         "test",
	Description:  "",
	Unrestricted: false,
	Secret:       "",
	ProjectID:    "53c2b94f63fb4f43a21b92d119ce549f",
	Roles: []applicationcredentials.Role{
		{
			ID:   "31f87923ae4a4d119aa0b85dcdbeed13",
			Name: "compute_viewer",
		},
	},
	AccessRules: []applicationcredentials.AccessRule{
		{
			ID:      "07d719df00f349ef8de77d542edf010c",
			Path:    "/v2.0/metrics",
			Method:  "GET",
			Service: "monitoring",
		},
	},
	ExpiresAt: nilTime,
	Links: map[string]any{
		"self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/f741662395b249c9b8acdebf1722c5ae",
	},
}

var ApplicationCredentialNoSecretResponse = applicationcredentials.ApplicationCredential{
	ID:           "c4859fb437df4b87a51a8f5adcfb0bc7",
	Name:         "test1",
	Description:  "",
	Unrestricted: false,
	Secret:       "generated_secret",
	ProjectID:    "53c2b94f63fb4f43a21b92d119ce549f",
	Roles: []applicationcredentials.Role{
		{
			ID:   "31f87923ae4a4d119aa0b85dcdbeed13",
			Name: "compute_viewer",
		},
	},
	ExpiresAt: nilTime,
	Links: map[string]any{
		"self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/c4859fb437df4b87a51a8f5adcfb0bc7",
	},
}

var ApplationCredentialExpiresAt = time.Date(2019, 3, 12, 12, 12, 12, 123456000, time.UTC)
var UnrestrictedApplicationCredential = applicationcredentials.ApplicationCredential{
	ID:           "6b8cc7647da64166a4a3cc0c88ebbabb",
	Name:         "test2",
	Description:  "",
	Unrestricted: true,
	Secret:       "",
	ProjectID:    "53c2b94f63fb4f43a21b92d119ce549f",
	Roles: []applicationcredentials.Role{
		{
			ID:   "31f87923ae4a4d119aa0b85dcdbeed13",
			Name: "compute_viewer",
		},
		{
			ID:   "4494bc5bea1a4105ad7fbba6a7eb9ef4",
			Name: "network_viewer",
		},
	},
	ExpiresAt: ApplationCredentialExpiresAt,
	Links: map[string]any{
		"self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/6b8cc7647da64166a4a3cc0c88ebbabb",
	},
}

var FirstApplicationCredential = applicationcredentials.ApplicationCredential{
	ID:           "c4859fb437df4b87a51a8f5adcfb0bc7",
	Name:         "test1",
	Description:  "",
	Unrestricted: false,
	Secret:       "",
	ProjectID:    "53c2b94f63fb4f43a21b92d119ce549f",
	Roles: []applicationcredentials.Role{
		{
			ID:   "31f87923ae4a4d119aa0b85dcdbeed13",
			Name: "compute_viewer",
		},
	},
	ExpiresAt: nilTime,
	Links: map[string]any{
		"self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/c4859fb437df4b87a51a8f5adcfb0bc7",
	},
}

var SecondApplicationCredential = applicationcredentials.ApplicationCredential{
	ID:           "6b8cc7647da64166a4a3cc0c88ebbabb",
	Name:         "test2",
	Description:  "",
	Unrestricted: true,
	Secret:       "",
	ProjectID:    "53c2b94f63fb4f43a21b92d119ce549f",
	Roles: []applicationcredentials.Role{
		{
			ID:   "31f87923ae4a4d119aa0b85dcdbeed13",
			Name: "compute_viewer",
		},
		{
			ID:   "4494bc5bea1a4105ad7fbba6a7eb9ef4",
			Name: "network_viewer",
		},
	},
	ExpiresAt: ApplationCredentialExpiresAt,
	Links: map[string]any{
		"self": "https://identity/v3/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/6b8cc7647da64166a4a3cc0c88ebbabb",
	},
}

var AccessRule = applicationcredentials.AccessRule{
	Path:    "/v2.0/metrics",
	ID:      "07d719df00f349ef8de77d542edf010c",
	Service: "monitoring",
	Method:  "GET",
}

var ExpectedAccessRulesSlice = []applicationcredentials.AccessRule{
	AccessRule,
}

// ExpectedApplicationCredentialsSlice is the slice of application credentials expected to be returned from ListOutput.
var ExpectedApplicationCredentialsSlice = []applicationcredentials.ApplicationCredential{FirstApplicationCredential, SecondApplicationCredential}

// HandleListApplicationCredentialsSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that responds with a list of two applicationcredentials.
func HandleListApplicationCredentialsSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/application_credentials", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, ListOutput)
	})
}

// HandleGetApplicationCredentialSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that responds with a single application credential.
func HandleGetApplicationCredentialSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/f741662395b249c9b8acdebf1722c5ae", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, GetOutput)
	})
}

// HandleCreateApplicationCredentialSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that tests application credential creation.
func HandleCreateApplicationCredentialSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/application_credentials", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, CreateResponse)
	})
}

// HandleCreateNoOptionsApplicationCredentialSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that tests application credential creation.
func HandleCreateNoSecretApplicationCredentialSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/application_credentials", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateNoSecretRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, CreateNoSecretResponse)
	})
}

func HandleCreateUnrestrictedApplicationCredentialSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/application_credentials", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateUnrestrictedRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, CreateUnrestrictedResponse)
	})
}

// HandleDeleteApplicationCredentialSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that tests application credential deletion.
func HandleDeleteApplicationCredentialSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/application_credentials/f741662395b249c9b8acdebf1722c5ae", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleListAccessRulesSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that responds with a list of two applicationcredentials.
func HandleListAccessRulesSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/access_rules", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, ListAccessRulesOutput)
	})
}

// HandleGetAccessRuleSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that responds with a single application credential.
func HandleGetAccessRuleSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/access_rules/07d719df00f349ef8de77d542edf010c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, GetAccessRuleOutput)
	})
}

// HandleDeleteAccessRuleSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that tests application credential deletion.
func HandleDeleteAccessRuleSuccessfully(t *testing.T, fakeServer th.FakeServer) {
	fakeServer.Mux.HandleFunc("/users/2844b2a08be147a08ef58317d6471f1f/access_rules/07d719df00f349ef8de77d542edf010c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}
