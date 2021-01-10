package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// ListAvailableOutput provides a single page of available Project results.
const ListAvailableOutput = `
{
  "projects": [
    {
      "description": "my first project",
      "domain_id": "11111",
      "enabled": true,
      "id": "abcde",
      "links": {
        "self": "http://localhost:5000/identity/v3/projects/abcde"
      },
      "name": "project 1",
      "parent_id": "11111"
    },
    {
      "description": "my second project",
      "domain_id": "22222",
      "enabled": true,
      "id": "bcdef",
      "links": {
        "self": "http://localhost:5000/identity/v3/projects/bcdef"
      },
      "name": "project 2",
      "parent_id": "22222"
    }
  ],
  "links": {
    "next": null,
    "previous": null,
    "self": "http://localhost:5000/identity/v3/users/foobar/projects"
  }
}
`

// ListOutput provides a single page of Project results.
const ListOutput = `
{
  "projects": [
    {
      "is_domain": false,
      "description": "The team that is red",
      "domain_id": "default",
      "enabled": true,
      "id": "1234",
      "name": "Red Team",
      "parent_id": null,
      "tags": ["Red", "Team"],
      "test": "old"
    },
    {
      "is_domain": false,
      "description": "The team that is blue",
      "domain_id": "default",
      "enabled": true,
      "id": "9876",
      "name": "Blue Team",
      "parent_id": null,
      "options": {
            "immutable": true
      }
    }
  ],
  "links": {
    "next": null,
    "previous": null
  }
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
  "project": {
		"is_domain": false,
		"description": "The team that is red",
		"domain_id": "default",
		"enabled": true,
		"id": "1234",
		"name": "Red Team",
		"parent_id": null,
		"tags": ["Red", "Team"],
		"test": "old"
	}
}
`

// CreateRequest provides the input to a Create request.
const CreateRequest = `
{
  "project": {
		"description": "The team that is red",
		"name": "Red Team",
		"tags": ["Red", "Team"],
		"test": "old"
  }
}
`

// UpdateRequest provides the input to an Update request.
const UpdateRequest = `
{
  "project": {
		"description": "The team that is bright red",
		"name": "Bright Red Team",
		"tags": ["Red"],
		"test": "new"
  }
}
`

// UpdateOutput provides an Update response.
const UpdateOutput = `
{
  "project": {
		"is_domain": false,
		"description": "The team that is bright red",
		"domain_id": "default",
		"enabled": true,
		"id": "1234",
		"name": "Bright Red Team",
		"parent_id": null,
		"tags": ["Red"],
		"test": "new"
	}
}
`

// FirstProject is a Project fixture.
var FirstProject = projects.Project{
	Description: "my first project",
	DomainID:    "11111",
	Enabled:     true,
	ID:          "abcde",
	Name:        "project 1",
	ParentID:    "11111",
	Extra: map[string]interface{}{
		"links": map[string]interface{}{"self": "http://localhost:5000/identity/v3/projects/abcde"},
	},
}

// SecondProject is a Project fixture.
var SecondProject = projects.Project{
	Description: "my second project",
	DomainID:    "22222",
	Enabled:     true,
	ID:          "bcdef",
	Name:        "project 2",
	ParentID:    "22222",
	Extra: map[string]interface{}{
		"links": map[string]interface{}{"self": "http://localhost:5000/identity/v3/projects/bcdef"},
	},
}

// RedTeam is a Project fixture.
var RedTeam = projects.Project{
	IsDomain:    false,
	Description: "The team that is red",
	DomainID:    "default",
	Enabled:     true,
	ID:          "1234",
	Name:        "Red Team",
	ParentID:    "",
	Tags:        []string{"Red", "Team"},
	Extra:       map[string]interface{}{"test": "old"},
}

// BlueTeam is a Project fixture.
var BlueTeam = projects.Project{
	IsDomain:    false,
	Description: "The team that is blue",
	DomainID:    "default",
	Enabled:     true,
	ID:          "9876",
	Name:        "Blue Team",
	ParentID:    "",
	Extra:       make(map[string]interface{}),
	Options: map[projects.Option]interface{}{
		projects.Immutable: true,
	},
}

// UpdatedRedTeam is a Project Fixture.
var UpdatedRedTeam = projects.Project{
	IsDomain:    false,
	Description: "The team that is bright red",
	DomainID:    "default",
	Enabled:     true,
	ID:          "1234",
	Name:        "Bright Red Team",
	ParentID:    "",
	Tags:        []string{"Red"},
	Extra:       map[string]interface{}{"test": "new"},
}

// ExpectedAvailableProjectsSlice is the slice of projects expected to be returned
// from ListAvailableOutput.
var ExpectedAvailableProjectsSlice = []projects.Project{FirstProject, SecondProject}

// ExpectedProjectSlice is the slice of projects expected to be returned from ListOutput.
var ExpectedProjectSlice = []projects.Project{RedTeam, BlueTeam}

// HandleListAvailableProjectsSuccessfully creates an HTTP handler at `/auth/projects`
// on the test handler mux that responds with a list of two tenants.
func HandleListAvailableProjectsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/auth/projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListAvailableOutput)
	})
}

// HandleListProjectsSuccessfully creates an HTTP handler at `/projects` on the
// test handler mux that responds with a list of two tenants.
func HandleListProjectsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleGetProjectSuccessfully creates an HTTP handler at `/projects` on the
// test handler mux that responds with a single project.
func HandleGetProjectSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/projects/1234", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleCreateProjectSuccessfully creates an HTTP handler at `/projects` on the
// test handler mux that tests project creation.
func HandleCreateProjectSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleDeleteProjectSuccessfully creates an HTTP handler at `/projects` on the
// test handler mux that tests project deletion.
func HandleDeleteProjectSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/projects/1234", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})
}

// HandleUpdateProjectSuccessfully creates an HTTP handler at `/projects` on the
// test handler mux that tests project updates.
func HandleUpdateProjectSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/projects/1234", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateOutput)
	})
}
