package services

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

func TestCreateSuccessful(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "X-Auth-Token", "1111")
		testhelper.TestJSONRequest(t, r, `{ "type": "compute" }`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `{
        "service": {
          "description": "Here's your service",
          "id": "1234",
          "name": "InscrutableOpenStackProjectName",
          "type": "compute"
        }
    }`)
	})

	client := gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{
			TokenID: "1111",
		},
		Endpoint: testhelper.Endpoint(),
	}

	result, err := Create(&client, "compute")
	if err != nil {
		t.Fatalf("Unexpected error from Create: %v", err)
	}

	if result.Description == nil || *result.Description != "Here's your service" {
		t.Errorf("Service description was unexpected [%s]", result.Description)
	}
	if result.ID != "1234" {
		t.Errorf("Service ID was unexpected [%s]", result.ID)
	}
	if result.Name != "InscrutableOpenStackProjectName" {
		t.Errorf("Service name was unexpected [%s]", result.Name)
	}
	if result.Type != "compute" {
		t.Errorf("Service type was unexpected [%s]", result.Type)
	}
}

func TestListSinglePage(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", "1111")

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
			{
				"links": {
					"next": null,
					"previous": null
				},
				"services": [
					{
						"description": "Service One",
						"id": "1234",
						"name": "service-one",
						"type": "identity"
					},
					{
						"description": "Service Two",
						"id": "9876",
						"name": "service-two",
						"type": "compute"
					}
				]
			}
		`)
	})

	client := gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{
			TokenID: "1111",
		},
		Endpoint: testhelper.Endpoint(),
	}

	result, err := List(&client, ListOpts{})
	if err != nil {
		t.Fatalf("Error listing services: %v", err)
	}

	if result.Pagination.Next != nil {
		t.Errorf("Unexpected next link: %s", result.Pagination.Next)
	}
	if result.Pagination.Previous != nil {
		t.Errorf("Unexpected previous link: %s", result.Pagination.Previous)
	}
	if len(result.Services) != 2 {
		t.Errorf("Unexpected number of services: %s", len(result.Services))
	}
	if result.Services[0].ID != "1234" {
		t.Errorf("Unexpected service: %#v", result.Services[0])
	}
	if result.Services[0].ID == "9876" {
		t.Errorf("Unexpected service: %#v", result.Services[1])
	}
}
