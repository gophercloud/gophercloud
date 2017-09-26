package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/groupresults"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/userresults"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// ListOutput provides a single page of Group results.
const ListOutput = `
{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/groups"
    },
    "groups": [
        {
            "domain_id": "default",
            "id": "2844b2a08be147a08ef58317d6471f1f",
            "description": "group for internal support users",
            "links": {
                "self": "http://example.com/identity/v3/groups/2844b2a08be147a08ef58317d6471f1f"
            },
            "name": "internal support",
            "extra": {
              "email": "support@localhost"
            }
        },
        {
            "domain_id": "1789d1",
            "id": "9fe1d3",
            "description": "group for support users",
            "links": {
                "self": "https://example.com/identity/v3/groups/9fe1d3"
            },
            "name": "support",
            "extra": {
              "email": "support@example.com"
            }
        }
    ]
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
    "group": {
        "domain_id": "1789d1",
        "id": "9fe1d3",
        "description": "group for support users",
        "links": {
            "self": "https://example.com/identity/v3/groups/9fe1d3"
        },
        "name": "support",
		"extra": {
              "email": "support@example.com"
        }
    }
}
`

// CreateRequest provides the input to a Create request.
const CreateRequest = `
{
    "group": {
        "domain_id": "1789d1",
        "name": "support",
        "description": "group for support users",
		"email": "support@example.com"
    }
}
`

// UpdateRequest provides the input to as Update request.
const UpdateRequest = `
{
    "group": {
        "description": "L2 Support Team"
    }
}
`

// UpdateOutput provides an update result.
const UpdateOutput = `
{
    "group": {
        "domain_id": "1789d1",
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/groups/9fe1d3"
        },
        "name": "support",
        "description": "L2 Support Team",
		"extra": {
              "email": "support@example.com"
        }
	}
}
`

// ListUsersOutput provides a ListUsers result.
const ListUsersOutput = `
{
    "users": [
        {
            "domain_id": "default",
            "enabled": true,
            "id": "ea167b",
            "links": {
                "self": "https://example.com/identity/v3/users/ea167b"
            },
            "name": "glance",
            "password_expires_at": null,
            "description": "some description",
            "extra": {
              "email": "glance@localhost"
            }
        },
        {
            "default_project_id": "263fd9",
            "domain_id": "1789d1",
            "enabled": true,
            "id": "a62db1",
            "links": {
                "self": "https://example.com/identity/v3/users/a62db1"
            },
            "name": "jsmith",
            "password_expires_at": "2016-11-06T15:32:17.000000",
            "email": "jsmith@example.com",
            "options": {
                "ignore_password_expiry": true,
                "multi_factor_auth_rules": [["password", "totp"], ["password", "custom-auth-method"]]
            }
        }
    ],
    "links": {
        "self": "https://example.com/identity/v3/groups/9fe1d3/users",
        "previous": null,
        "next": null
    }
}
`

// FirstGroup is the first group in the List request.
var FirstGroup = groupresults.Group{
	DomainID: "default",
	ID:       "2844b2a08be147a08ef58317d6471f1f",
	Links: map[string]interface{}{
		"self": "http://example.com/identity/v3/groups/2844b2a08be147a08ef58317d6471f1f",
	},
	Name:        "internal support",
	Description: "group for internal support users",
	Extra: map[string]interface{}{
		"email": "support@localhost",
	},
}

// SecondGroup is the second group in the List request.
var SecondGroup = groupresults.Group{
	DomainID: "1789d1",
	ID:       "9fe1d3",
	Links: map[string]interface{}{
		"self": "https://example.com/identity/v3/groups/9fe1d3",
	},
	Name:        "support",
	Description: "group for support users",
	Extra: map[string]interface{}{
		"email": "support@example.com",
	},
}

// SecondGroupUpdated is how SecondGroup should look after an Update.
var SecondGroupUpdated = groupresults.Group{
	DomainID: "1789d1",
	ID:       "9fe1d3",
	Links: map[string]interface{}{
		"self": "https://example.com/identity/v3/groups/9fe1d3",
	},
	Name:        "support",
	Description: "L2 Support Team",
	Extra: map[string]interface{}{
		"email": "support@example.com",
	},
}

// ExpectedGroupsSlice is the slice of groups expected to be returned from ListOutput.
var ExpectedGroupsSlice = []groupresults.Group{FirstGroup, SecondGroup}

var nilTime time.Time
var FirstUser = userresults.User{
	DomainID: "default",
	Enabled:  true,
	ID:       "ea167b",
	Links: map[string]interface{}{
		"self": "https://example.com/identity/v3/users/ea167b",
	},
	Name:              "glance",
	PasswordExpiresAt: nilTime,
	Description:       "some description",
	Extra: map[string]interface{}{
		"email": "glance@localhost",
	},
}

// SecondUser is the second user in the List request.
var SecondUserPasswordExpiresAt, _ = time.Parse(gophercloud.RFC3339MilliNoZ, "2016-11-06T15:32:17.000000")
var SecondUser = userresults.User{
	DefaultProjectID: "263fd9",
	DomainID:         "1789d1",
	Enabled:          true,
	ID:               "a62db1",
	Links: map[string]interface{}{
		"self": "https://example.com/identity/v3/users/a62db1",
	},
	Name:              "jsmith",
	PasswordExpiresAt: SecondUserPasswordExpiresAt,
	Extra: map[string]interface{}{
		"email": "jsmith@example.com",
	},
	Options: map[string]interface{}{
		"ignore_password_expiry": true,
		"multi_factor_auth_rules": []interface{}{
			[]string{"password", "totp"},
			[]string{"password", "custom-auth-method"},
		},
	},
}

var ExpectedUsersSlice = []userresults.User{FirstUser, SecondUser}

// HandleListGroupsSuccessfully creates an HTTP handler at `/groups` on the
// test handler mux that responds with a list of two groups.
func HandleListGroupsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetGroupSuccessfully creates an HTTP handler at `/groups` on the
// test handler mux that responds with a single group.
func HandleGetGroupSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/groups/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleCreateGroupSuccessfully creates an HTTP handler at `/groups` on the
// test handler mux that tests group creation.
func HandleCreateGroupSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleUpdateGroupSuccessfully creates an HTTP handler at `/groups` on the
// test handler mux that tests group update.
func HandleUpdateGroupSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/groups/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateOutput)
	})
}

// HandleDeleteGroupSuccessfully creates an HTTP handler at `/groups` on the
// test handler mux that tests group deletion.
func HandleDeleteGroupSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/groups/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleListGroupUsersSuccessfully creates an HTTP handler at /groups/{groupID}/users
// on the test handler mux that respons wit a list of two users
func HandleListGroupUsersSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/groups/9fe1d3/users", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListUsersOutput)
	})
}
