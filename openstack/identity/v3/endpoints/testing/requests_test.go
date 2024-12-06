package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateSuccessful(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
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
		fmt.Fprint(w, `
      {
        "endpoint": {
          "id": "12",
          "interface": "public",
		  "enabled": true,
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

	actual, err := endpoints.Create(context.TODO(), client.ServiceClient(), endpoints.CreateOpts{
		Availability: gophercloud.AvailabilityPublic,
		Name:         "the-endiest-of-points",
		Region:       "underground",
		URL:          "https://1.2.3.4:9000/",
		ServiceID:    "asdfasdfasdfasdf",
	}).Extract()
	th.AssertNoErr(t, err)

	expected := &endpoints.Endpoint{
		ID:           "12",
		Availability: gophercloud.AvailabilityPublic,
		Enabled:      true,
		Name:         "the-endiest-of-points",
		Region:       "underground",
		ServiceID:    "asdfasdfasdfasdf",
		URL:          "https://1.2.3.4:9000/",
	}

	th.AssertDeepEquals(t, expected, actual)
}

func TestListEndpoints(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, `
			{
				"endpoints": [
					{
						"id": "12",
						"interface": "public",
						"enabled": true,
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
						"enabled": false,
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

	count := 0
	err := endpoints.List(client.ServiceClient(), endpoints.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := endpoints.ExtractEndpoints(page)
		if err != nil {
			t.Errorf("Failed to extract endpoints: %v", err)
			return false, err
		}

		expected := []endpoints.Endpoint{
			{
				ID:           "12",
				Availability: gophercloud.AvailabilityPublic,
				Enabled:      true,
				Name:         "the-endiest-of-points",
				Region:       "underground",
				ServiceID:    "asdfasdfasdfasdf",
				URL:          "https://1.2.3.4:9000/",
			},
			{
				ID:           "13",
				Availability: gophercloud.AvailabilityInternal,
				Enabled:      false,
				Name:         "shhhh",
				Region:       "underground",
				ServiceID:    "asdfasdfasdfasdf",
				URL:          "https://1.2.3.4:9001/",
			},
		}
		th.AssertDeepEquals(t, expected, actual)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestUpdateEndpoint(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/endpoints/12", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, `
		{
	    "endpoint": {
	      "name": "renamed",
				"region": "somewhere-else"
	    }
		}
	`)

		fmt.Fprint(w, `
		{
			"endpoint": {
				"id": "12",
				"interface": "public",
				"enabled": true,
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

	actual, err := endpoints.Update(context.TODO(), client.ServiceClient(), "12", endpoints.UpdateOpts{
		Name:   "renamed",
		Region: "somewhere-else",
	}).Extract()
	if err != nil {
		t.Fatalf("Unexpected error from Update: %v", err)
	}

	expected := &endpoints.Endpoint{
		ID:           "12",
		Availability: gophercloud.AvailabilityPublic,
		Enabled:      true,
		Name:         "renamed",
		Region:       "somewhere-else",
		ServiceID:    "asdfasdfasdfasdf",
		URL:          "https://1.2.3.4:9000/",
	}
	th.AssertDeepEquals(t, expected, actual)
}

func TestDeleteEndpoint(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/endpoints/34", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	res := endpoints.Delete(context.TODO(), client.ServiceClient(), "34")
	th.AssertNoErr(t, res.Err)
}
