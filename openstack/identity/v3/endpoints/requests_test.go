package endpoints

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/testhelper"
)

const tokenID = "abcabcabcabc"

func serviceClient() *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenID},
		Endpoint: testhelper.Endpoint(),
	}
}

func TestCreateSuccessful(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "POST")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenID)
		testhelper.TestJSONRequest(t, r, `
      {
        "endpoint": {
          "interface": "public",
          "name": "the-endiest-of-points",
          "region": "underground",
          "url": "https://1.2.3.4:9000/",
          "service_id": "asdfasdfasdfasdf"
        }
      }
    `)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, `
      {
        "endpoint": {
          "id": "12",
          "interface": "public",
          "links": {
            "self": "https://localhost:5000/v3/endpoints/12"
          },
          "name": "the-endiest-of-points",
          "region": "underground",
          "service_id": "asdfasdfasdfasdf",
          "url": "https://1.2.3.4:9000/"
        }
      }
    `)
	})

	client := serviceClient()

	result, err := Create(client, EndpointOpts{
		Interface: InterfacePublic,
		Name:      "the-endiest-of-points",
		Region:    "underground",
		URL:       "https://1.2.3.4:9000/",
		ServiceID: "asdfasdfasdfasdf",
	})
	if err != nil {
		t.Fatalf("Unable to create an endpoint: %v", err)
	}

	expected := &Endpoint{
		ID:        "12",
		Interface: InterfacePublic,
		Name:      "the-endiest-of-points",
		Region:    "underground",
		ServiceID: "asdfasdfasdfasdf",
		URL:       "https://1.2.3.4:9000/",
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %#v, was %#v", expected, result)
	}
}

func TestListEndpoints(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenID)

		fmt.Fprintf(w, `
			[
				{
					"id": "12",
					"interface": "public",
					"links": {
						"self": "https://localhost:5000/v3/endpoints/12"
					},
					"name": "the-endiest-of-points",
					"region": "underground",
					"service_id": "asdfasdfasdfasdf",
					"url": "https://1.2.3.4:9000/"
				},
				{
					"id": "13",
					"interface": "internal",
					"links": {
						"self": "https://localhost:5000/v3/endpoints/13"
					},
					"name": "shhhh",
					"region": "underground",
					"service_id": "asdfasdfasdfasdf",
					"url": "https://1.2.3.4:9001/"
				}
			]
		`)
	})

	client := serviceClient()

	actual, err := List(client, ListOpts{})
	if err != nil {
		t.Fatalf("Unexpected error listing endpoints: %v", err)
	}

	expected := &EndpointList{
		Endpoints: []Endpoint{
			Endpoint{
				ID:        "12",
				Interface: InterfacePublic,
				Name:      "the-endiest-of-points",
				Region:    "underground",
				ServiceID: "asdfasdfasdfasdf",
				URL:       "https://1.2.3.4:9000/",
			},
			Endpoint{
				ID:        "13",
				Interface: InterfaceInternal,
				Name:      "shhhh",
				Region:    "underground",
				ServiceID: "asdfasdfasdfasdf",
				URL:       "https://1.2.3.4:9001/",
			},
		},
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func TestUpdateEndpoint(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/endpoints/12", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "GET")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenID)
		testhelper.TestJSONRequest(t, r, `
		{
	    "endpoint": {
	      "name": "renamed",
				"region": "somewhere-else"
	    }
		}
	`)

		fmt.Fprintf(w, `
		{
			"endpoint": {
				"id": "12",
				"interface": "public",
				"links": {
					"self": "https://localhost:5000/v3/endpoints/12"
				},
				"name": "renamed",
				"region": "somewhere-else",
				"service_id": "asdfasdfasdfasdf",
				"url": "https://1.2.3.4:9000/"
			}
		}
	`)
	})

	client := serviceClient()
	actual, err := Update(client, "12", EndpointOpts{
		Name:   "renamed",
		Region: "somewhere-else",
	})
	if err != nil {
		t.Fatalf("Unexpected error from Update: %v", err)
	}

	expected := Endpoint{
		ID:        "12",
		Interface: InterfacePublic,
		Name:      "renamed",
		Region:    "somewhere-else",
		ServiceID: "asdfasdfasdfasdf",
		URL:       "https://1.2.3.4:9000/",
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %#v, was %#v", expected, actual)
	}
}

func TestDeleteEndpoint(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/endpoints/34", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "DELETE")
		testhelper.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	client := serviceClient()

	err := Delete(client, "34")
	if err != nil {
		t.Fatalf("Unexpected error from Delete: %v", err)
	}
}
