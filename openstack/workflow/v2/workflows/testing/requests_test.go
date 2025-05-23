package testing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/workflow/v2/workflows"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateWorkflow(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	definition := `---
version: '2.0'

workflow_echo:
	description: Simple workflow example
	type: direct
	input:
		- msg

	tasks:
		test:
        	action: std.echo output="<% $.msg %>"`

	fakeServer.Mux.HandleFunc("/workflows", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestHeader(t, r, "Content-Type", "text/plain")
		th.TestFormValues(t, r, map[string]string{
			"namespace": "some-namespace",
			"scope":     "private",
		})
		th.TestBody(t, r, definition)

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")

		fmt.Fprint(w, `{
			"workflows": [
				{
					"created_at": "2018-09-12 15:48:17",
					"definition": "---\nversion: '2.0'\n\nworkflow_echo:\n  description: Simple workflow example\n  type: direct\n\n  input:\n    - msg\n\n  tasks:\n    test:\n      action: std.echo output=\"<% $.msg %>\"",
					"id": "604a3a1e-94e3-4066-a34a-aa56873ef236",
					"input": "msg",
					"name": "workflow_echo",
					"namespace": "some-namespace",
					"project_id": "778c0f25df0d492a9a868ee9e2fbb513",
					"scope": "private",
					"tags": [],
					"updated_at": "2018-09-12 15:48:17"
				}
			]
		}`)
	})

	opts := &workflows.CreateOpts{
		Namespace:  "some-namespace",
		Scope:      "private",
		Definition: strings.NewReader(definition),
	}

	actual, err := workflows.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	if err != nil {
		t.Fatalf("Unable to create workflow: %v", err)
	}

	updated := time.Date(2018, time.September, 12, 15, 48, 17, 0, time.UTC)
	expected := []workflows.Workflow{
		{
			ID:         "604a3a1e-94e3-4066-a34a-aa56873ef236",
			Definition: "---\nversion: '2.0'\n\nworkflow_echo:\n  description: Simple workflow example\n  type: direct\n\n  input:\n    - msg\n\n  tasks:\n    test:\n      action: std.echo output=\"<% $.msg %>\"",
			Name:       "workflow_echo",
			Namespace:  "some-namespace",
			Input:      "msg",
			ProjectID:  "778c0f25df0d492a9a868ee9e2fbb513",
			Scope:      "private",
			Tags:       []string{},
			CreatedAt:  time.Date(2018, time.September, 12, 15, 48, 17, 0, time.UTC),
			UpdatedAt:  &updated,
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestDeleteWorkflow(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/workflows/604a3a1e-94e3-4066-a34a-aa56873ef236", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusAccepted)
	})

	res := workflows.Delete(context.TODO(), client.ServiceClient(fakeServer), "604a3a1e-94e3-4066-a34a-aa56873ef236")
	th.AssertNoErr(t, res.Err)
}

func TestGetWorkflow(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fakeServer.Mux.HandleFunc("/workflows/1", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"created_at": "2018-09-12 15:48:17",
				"definition": "---\nversion: '2.0'\n\nworkflow_echo:\n  description: Simple workflow example\n  type: direct\n\n  input:\n    - msg\n\n  tasks:\n    test:\n      action: std.echo output=\"<% $.msg %>\"",
				"id": "604a3a1e-94e3-4066-a34a-aa56873ef236",
				"input": "msg",
				"name": "workflow_echo",
				"namespace": "some-namespace",
				"project_id": "778c0f25df0d492a9a868ee9e2fbb513",
				"scope": "private",
				"tags": [],
				"updated_at": "2018-09-12 15:48:17"
			}
		`)
	})
	actual, err := workflows.Get(context.TODO(), client.ServiceClient(fakeServer), "1").Extract()
	if err != nil {
		t.Fatalf("Unable to get workflow: %v", err)
	}

	updated := time.Date(2018, time.September, 12, 15, 48, 17, 0, time.UTC)
	expected := &workflows.Workflow{
		ID:         "604a3a1e-94e3-4066-a34a-aa56873ef236",
		Definition: "---\nversion: '2.0'\n\nworkflow_echo:\n  description: Simple workflow example\n  type: direct\n\n  input:\n    - msg\n\n  tasks:\n    test:\n      action: std.echo output=\"<% $.msg %>\"",
		Name:       "workflow_echo",
		Namespace:  "some-namespace",
		Input:      "msg",
		ProjectID:  "778c0f25df0d492a9a868ee9e2fbb513",
		Scope:      "private",
		Tags:       []string{},
		CreatedAt:  time.Date(2018, time.September, 12, 15, 48, 17, 0, time.UTC),
		UpdatedAt:  &updated,
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, but was %#v", expected, actual)
	}
}

func TestListWorkflows(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fakeServer.Mux.HandleFunc("/workflows", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		w.Header().Add("Content-Type", "application/json")
		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse request form %v", err)
		}
		marker := r.Form.Get("marker")
		switch marker {
		case "":
			fmt.Fprintf(w, `{
				"next": "%s/workflows?marker=604a3a1e-94e3-4066-a34a-aa56873ef236",
				"workflows": [
					{
						"created_at": "2018-09-12 15:48:17",
						"definition": "---\nversion: '2.0'\n\nworkflow_echo:\n  description: Simple workflow example\n  type: direct\n\n  input:\n    - msg\n\n  tasks:\n    test:\n      action: std.echo output=\"<%% $.msg %%>\"",
						"id": "604a3a1e-94e3-4066-a34a-aa56873ef236",
						"input": "msg",
						"name": "workflow_echo",
						"namespace": "some-namespace",
						"project_id": "778c0f25df0d492a9a868ee9e2fbb513",
						"scope": "private",
						"tags": [],
						"updated_at": "2018-09-12 15:48:17"
					}
				]
			}`, fakeServer.Server.URL)
		case "604a3a1e-94e3-4066-a34a-aa56873ef236":
			fmt.Fprint(w, `{ "workflows": [] }`)
		default:
			t.Fatalf("Unexpected marker: [%s]", marker)
		}
	})
	pages := 0
	// Get all workflows
	err := workflows.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++
		actual, err := workflows.ExtractWorkflows(page)
		if err != nil {
			return false, err
		}

		updated := time.Date(2018, time.September, 12, 15, 48, 17, 0, time.UTC)
		expected := []workflows.Workflow{
			{
				ID:         "604a3a1e-94e3-4066-a34a-aa56873ef236",
				Definition: "---\nversion: '2.0'\n\nworkflow_echo:\n  description: Simple workflow example\n  type: direct\n\n  input:\n    - msg\n\n  tasks:\n    test:\n      action: std.echo output=\"<% $.msg %>\"",
				Name:       "workflow_echo",
				Namespace:  "some-namespace",
				Input:      "msg",
				ProjectID:  "778c0f25df0d492a9a868ee9e2fbb513",
				Scope:      "private",
				Tags:       []string{},
				CreatedAt:  time.Date(2018, time.September, 12, 15, 48, 17, 0, time.UTC),
				UpdatedAt:  &updated,
			},
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %#v, but was %#v", expected, actual)
		}
		return true, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if pages != 1 {
		t.Errorf("Expected one page, got %d", pages)
	}
}

func TestToWorkflowListQuery(t *testing.T) {
	for expected, opts := range map[string]*workflows.ListOpts{
		newValue("tags", `tag1,tag2`): {
			Tags: []string{"tag1", "tag2"},
		},
		newValue("name", `neq:invalid_name`): {
			Name: &workflows.ListFilter{
				Filter: workflows.FilterNEQ,
				Value:  "invalid_name",
			},
		},
		newValue("created_at", `gt:2018-01-01 00:00:00`): {
			CreatedAt: &workflows.ListDateFilter{
				Filter: workflows.FilterGT,
				Value:  time.Date(2018, time.January, 1, 0, 0, 0, 0, time.UTC),
			},
		},
	} {
		actual, _ := opts.ToWorkflowListQuery()
		th.AssertEquals(t, expected, actual)
	}
}
func newValue(param, value string) string {
	v := url.Values{}
	v.Add(param, value)
	return "?" + v.Encode()
}
