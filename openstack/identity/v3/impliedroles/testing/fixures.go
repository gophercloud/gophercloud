package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/impliedroles"
	"github.com/gophercloud/gophercloud/testhelper"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

const ListOutput = `
{
    "role_inferences": [{
        "prior_role": {
            "id": "3b2800ee2ed44df88b76f677c8c760fb",
            "links": {
                "self": "http://example.com/v3/roles/3b2800ee2ed44df88b76f677c8c760fb"
            },
            "name": "member"
        },
        "implies": [{
            "id": "66f7869ae8134dbd89279fcdb43a5993",
            "links": {
                "self": "http://example.com/v3/roles/66f7869ae8134dbd89279fcdb43a5993"
            },
            "name": "reader"
        }]
    }, {
        "prior_role": {
            "id": "b385b97c988f4a649eecbb5cdd52b7e1",
            "links": {
                "self": "http://example.com/v3/roles/b385b97c988f4a649eecbb5cdd52b7e1"
            },
            "name": "admin"
        },
        "implies": [{
            "id": "3b2800ee2ed44df88b76f677c8c760fb",
            "links": {
                "self": "http://example.com/v3/roles/3b2800ee2ed44df88b76f677c8c760fb"
            },
            "name": "member"
        }]
    }]
}
`
const CreateResponse = `
{
	"role_inference": {
		"prior_role": {
			"id": "b385b97c988f4a649eecbb5cdd52b7e1", 
			"links": {
				"self": "http://example.com/v3/roles/b385b97c988f4a649eecbb5cdd52b7e1"
				}, 
			"name": "admin"
		}, 
		"implies": {
			"id": "ddb5331895e348b0ab78cf0db18e8b78", 
			"links": {
				"self": "http://example.com/v3/roles/ddb5331895e348b0ab78cf0db18e8b78"
			}, 
			"name": "test"
		}
	},
	"links": {
		"self": "http://example.com/v3/roles/b385b97c988f4a649eecbb5cdd52b7e1/implies/ddb5331895e348b0ab78cf0db18e8b78"
	}
}
`

var FirstImpliedRole = impliedroles.ImpliedRole{
	PriorRole: struct {
		ID    string                 "json:\"id\""
		Name  string                 "json:\"name\""
		Links map[string]interface{} "json:\"links\""
	}{
		ID:   "3b2800ee2ed44df88b76f677c8c760fb",
		Name: "member",
		Links: map[string]interface{}{
			"self": "http://example.com/v3/roles/3b2800ee2ed44df88b76f677c8c760fb",
		},
	},
	Implies: []struct {
		ID    string                 "json:\"id\""
		Name  string                 "json:\"name\""
		Links map[string]interface{} "json:\"links\""
	}{
		{
			ID:   "66f7869ae8134dbd89279fcdb43a5993",
			Name: "reader",
			Links: map[string]interface{}{
				"self": "http://example.com/v3/roles/66f7869ae8134dbd89279fcdb43a5993",
			},
		},
	},
}

var SecondImpliedRole = impliedroles.ImpliedRole{
	PriorRole: struct {
		ID    string                 "json:\"id\""
		Name  string                 "json:\"name\""
		Links map[string]interface{} "json:\"links\""
	}{
		ID:   "b385b97c988f4a649eecbb5cdd52b7e1",
		Name: "admin",
		Links: map[string]interface{}{
			"self": "http://example.com/v3/roles/b385b97c988f4a649eecbb5cdd52b7e1",
		},
	},
	Implies: []struct {
		ID    string                 "json:\"id\""
		Name  string                 "json:\"name\""
		Links map[string]interface{} "json:\"links\""
	}{
		{
			ID:   "3b2800ee2ed44df88b76f677c8c760fb",
			Name: "member",
			Links: map[string]interface{}{
				"self": "http://example.com/v3/roles/3b2800ee2ed44df88b76f677c8c760fb",
			},
		},
	},
}

var CreateImpliedRole = impliedroles.CreateImpliedRole{
	PriorRole: struct {
		ID    string                 "json:\"id\""
		Name  string                 "json:\"name\""
		Links map[string]interface{} "json:\"links\""
	}{
		ID:   "b385b97c988f4a649eecbb5cdd52b7e1",
		Name: "admin",
		Links: map[string]interface{}{
			"self": "http://example.com/v3/roles/b385b97c988f4a649eecbb5cdd52b7e1",
		},
	},
	Implies: struct {
		ID    string                 "json:\"id\""
		Name  string                 "json:\"name\""
		Links map[string]interface{} "json:\"links\""
	}{
		ID:   "ddb5331895e348b0ab78cf0db18e8b78",
		Name: "test",
		Links: map[string]interface{}{
			"self": "http://example.com/v3/roles/ddb5331895e348b0ab78cf0db18e8b78",
		},
	},
}

// HandleCreateImpliedRolesSuccessfully creates an HTTP handler at `/role_inferences` on the
// test handler mux that responds with a list of two roles.
func HandleListImpliedRolesSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/role_inferences", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

var ExpectedImpliedRoleSlice = []impliedroles.ImpliedRole{FirstImpliedRole, SecondImpliedRole}

// HandleCreateImpliedRolesSuccessfully Creates an HTTP handler at `/roles/{piror_role_id}/implies/{implied_role_id}`
// on the test handler mux that responds with response from the create implies role
func HandleCreateImpliedRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles/b385b97c988f4a649eecbb5cdd52b7e1/implies/ddb5331895e348b0ab78cf0db18e8b78", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, err := fmt.Fprintf(w, CreateResponse)
		testhelper.AssertNoErr(t, err)
	})
}

// HandleDeleteImpliedRole Create an HTTP handler at '/roles/{prior_role_id}/implies/{implied_role_id}`
// on test handler Mux that responds with the response from the delete implies role
func HandleDeleteImpliedRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles/b385b97c988f4a649eecbb5cdd52b7e1/implies/ddb5331895e348b0ab78cf0db18e8b78", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}
