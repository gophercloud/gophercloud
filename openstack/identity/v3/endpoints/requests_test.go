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
		Availability: gophercloud.AvailabilityPublic,
		Name:         "the-endiest-of-points",
		Region:       "underground",
		URL:          "https://1.2.3.4:9000/",
		ServiceID:    "asdfasdfasdfasdf",
	})
	if err != nil {
		t.Fatalf("Unable to create an endpoint: %v", err)
	}

	expected := &Endpoint{
		ID:           "12",
		Availability: gophercloud.AvailabilityPublic,
		Name:         "the-endiest-of-points",
		Region:       "underground",
		ServiceID:    "asdfasdfasdfasdf",
		URL:          "https://1.2.3.4:9000/",
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
			{
				"endpoints": [
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
				],
				"links": {
					"next": null,
					"previous": null
				}
			}
		`)
	})

	client := serviceClient()

	count := 0
	List(client, ListOpts{}).EachPage(func(page gophercloud.Page) bool {
		count++
		actual, err := ExtractEndpoints(page)
		if err != nil {
			t.Errorf("Failed to extract endpoints: %v", err)
			return false
		}

		expected := []Endpoint{
			Endpoint{
				ID:           "12",
				Availability: gophercloud.AvailabilityPublic,
				Name:         "the-endiest-of-points",
				Region:       "underground",
				ServiceID:    "asdfasdfasdfasdf",
				URL:          "https://1.2.3.4:9000/",
			},
			Endpoint{
				ID:           "13",
				Availability: gophercloud.AvailabilityInternal,
				Name:         "shhhh",
				Region:       "underground",
				ServiceID:    "asdfasdfasdfasdf",
				URL:          "https://1.2.3.4:9001/",
			},
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %#v, got %#v", expected, actual)
		}

		return true
	})
	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestUpdateEndpoint(t *testing.T) {
	testhelper.SetupHTTP()
	defer testhelper.TeardownHTTP()

	testhelper.Mux.HandleFunc("/endpoints/12", func(w http.ResponseWriter, r *http.Request) {
		testhelper.TestMethod(t, r, "PATCH")
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

	expected := &Endpoint{
		ID:           "12",
		Availability: gophercloud.AvailabilityPublic,
		Name:         "renamed",
		Region:       "somewhere-else",
		ServiceID:    "asdfasdfasdfasdf",
		URL:          "https://1.2.3.4:9000/",
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
