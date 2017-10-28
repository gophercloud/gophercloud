package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

// ListOutput provides a single page of Role results.
const ListOutput = `
{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/roles"
    },
    "roles": [
        {
            "domain_id": "default",
            "id": "2844b2a08be147a08ef58317d6471f1f",
            "links": {
                "self": "http://example.com/identity/v3/roles/2844b2a08be147a08ef58317d6471f1f"
            },
            "name": "admin-read-only"
        },
        {
            "domain_id": "1789d1",
            "id": "9fe1d3",
            "links": {
                "self": "https://example.com/identity/v3/roles/9fe1d3"
            },
            "name": "support",
            "extra": {
                "description": "read-only support role"
            }
        }
    ]
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
    "role": {
        "domain_id": "1789d1",
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/roles/9fe1d3"
        },
        "name": "support",
        "extra": {
            "description": "read-only support role"
        }
    }
}
`

// CreateRequest provides the input to a Create request.
const CreateRequest = `
{
    "role": {
        "domain_id": "1789d1",
        "name": "support",
        "description": "read-only support role"
    }
}
`

// FirstRole is the first role in the List request.
var FirstRole = roles.Role{
	DomainID: "default",
	ID:       "2844b2a08be147a08ef58317d6471f1f",
	Links: map[string]interface{}{
		"self": "http://example.com/identity/v3/roles/2844b2a08be147a08ef58317d6471f1f",
	},
	Name:  "admin-read-only",
	Extra: map[string]interface{}{},
}

// SecondRole is the second role in the List request.
var SecondRole = roles.Role{
	DomainID: "1789d1",
	ID:       "9fe1d3",
	Links: map[string]interface{}{
		"self": "https://example.com/identity/v3/roles/9fe1d3",
	},
	Name: "support",
	Extra: map[string]interface{}{
		"description": "read-only support role",
	},
}

// ExpectedRolesSlice is the slice of roles expected to be returned from ListOutput.
var ExpectedRolesSlice = []roles.Role{FirstRole, SecondRole}

// HandleListRolesSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that responds with a list of two roles.
func HandleListRolesSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetRoleSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that responds with a single role.
func HandleGetRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleCreateRoleSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that tests role creation.
func HandleCreateRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleDeleteRoleSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that tests role deletion.
func HandleDeleteRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleAssignSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/projects/{project_id}/users/{user_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/projects/{project_id}/groups/{group_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/domains/{domain_id}/users/{user_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/domains/{domain_id}/groups/{group_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}

func HandleUnassignSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/projects/{project_id}/users/{user_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/projects/{project_id}/groups/{group_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/domains/{domain_id}/users/{user_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	th.Mux.HandleFunc("/domains/{domain_id}/groups/{group_id}/roles/{role_id}", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})
}
