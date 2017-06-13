package testing

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// ListOutput provides a single page of User results.
const ListOutput = `
{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/users"
    },
    "users": [
        {
            "domain_id": "default",
            "enabled": true,
            "id": "2844b2a08be147a08ef58317d6471f1f",
            "links": {
                "self": "http://example.com/identity/v3/users/2844b2a08be147a08ef58317d6471f1f"
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
            "id": "9fe1d3",
            "links": {
                "self": "https://example.com/identity/v3/users/9fe1d3"
            },
            "name": "jsmith",
            "password_expires_at": "2016-11-06T15:32:17.000000",
            "email": "jsmith@example.com",
            "options": {
                "ignore_password_expiry": true
            }
        }
    ]
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
    "user": {
        "default_project_id": "263fd9",
        "domain_id": "1789d1",
        "enabled": true,
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/users/9fe1d3"
        },
        "name": "jsmith",
        "password_expires_at": "2016-11-06T15:32:17.000000",
        "email": "jsmith@example.com",
        "options": {
            "ignore_password_expiry": true
        }
    }
}
`

// FirstUser is the first user in the List request.
var nilTime time.Time
var FirstUser = users.User{
	DomainID: "default",
	Enabled:  true,
	ID:       "2844b2a08be147a08ef58317d6471f1f",
	Links: map[string]interface{}{
		"self": "http://example.com/identity/v3/users/2844b2a08be147a08ef58317d6471f1f",
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
var SecondUser = users.User{
	DefaultProjectID: "263fd9",
	DomainID:         "1789d1",
	Enabled:          true,
	ID:               "9fe1d3",
	Links: map[string]interface{}{
		"self": "https://example.com/identity/v3/users/9fe1d3",
	},
	Name:              "jsmith",
	PasswordExpiresAt: SecondUserPasswordExpiresAt,
	Extra: map[string]interface{}{
		"email": "jsmith@example.com",
	},
	Options: map[string]interface{}{
		"ignore_password_expiry": true,
	},
}

// ExpectedUsersSlice is the slice of users expected to be returned from ListOutput.
var ExpectedUsersSlice = []users.User{FirstUser, SecondUser}

// HandleListUsersSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that responds with a list of two users.
func HandleListUsersSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetUserSuccessfully creates an HTTP handler at `/users` on the
// test handler mux that responds with a single user.
func HandleGetUserSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/users/9fe1d3", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}
