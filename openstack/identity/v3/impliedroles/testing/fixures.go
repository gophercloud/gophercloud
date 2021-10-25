package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/impliedroles"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

const ListOutput = `
{
    "role_inference": {
        "prior_role": {
            "id": "42c764f0c19146728dbfe73a49cc35c3",
            "links": {
                "self": "http://example.com/identity/v3/roles/42c764f0c19146728dbfe73a49cc35c3"
            },
            "name": "prior role name"
        },
        "implies": [
            {
                "id": "066fbfc8b3e54fb68784c9e7e92ab8d7",
                "links": {
                    "self": "http://example.com/identity/v3/roles/066fbfc8b3e54fb68784c9e7e92ab8d7"
                },
                "name": "implied role1 name"
            },
            {
                "id": "32a0df1cc22848aca3986adae9e0b9a0",
                "links": {
                    "self": "http://example.com/identity/v3/roles/32a0df1cc22848aca3986adae9e0b9a0"
                },
                "name": "implied role2 name"
            }
        ]
    },
    "links" : {
        "self": "http://example.com/identity/v3/roles/42c764f0c19146728dbfe73a49cc35c3/implies"
    }
}
`

var GetImpliedRole = impliedroles.GetImpliedRole{
	RoleInference: impliedroles.ImpliedRole{
		PriorRole: impliedroles.Role{
			ID:   "42c764f0c19146728dbfe73a49cc35c3",
			Name: "prior role name",
			Links: map[string]interface{}{
				"self": "http://example.com/identity/v3/roles/42c764f0c19146728dbfe73a49cc35c3",
			},
		},
		Implies: []impliedroles.Role{
			{
				ID:   "066fbfc8b3e54fb68784c9e7e92ab8d7",
				Name: "implied role1 name",
				Links: map[string]interface{}{
					"self": "http://example.com/identity/v3/roles/066fbfc8b3e54fb68784c9e7e92ab8d7",
				},
			}, {
				ID:   "32a0df1cc22848aca3986adae9e0b9a0",
				Name: "implied role2 name",
				Links: map[string]interface{}{
					"self": "http://example.com/identity/v3/roles/32a0df1cc22848aca3986adae9e0b9a0",
				},
			},
		},
	},
	Links: map[string]interface{}{
		"self": "http://example.com/identity/v3/roles/42c764f0c19146728dbfe73a49cc35c3/implies",
	},
}

// HandleCreateImpliedRolesSuccessfully creates an HTTP handler at `/role_inferences` on the
// test handler mux that responds with a list of two roles.
func HandleListImpliedRolesSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles/42c764f0c19146728dbfe73a49cc35c3/implies", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}
