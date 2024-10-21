package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projectendpoints"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreateSuccessful(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/OS-EP-FILTER/projects/project-id/endpoints/endpoint-id", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	err := projectendpoints.Create(context.TODO(), client.ServiceClient(), "project-id", "endpoint-id").Err
	th.AssertNoErr(t, err)
}

func TestListEndpoints(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/OS-EP-FILTER/projects/project-id/endpoints", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
		{
			"endpoints": [
				{
					"id": "6fedc0",
					"interface": "public",
					"url": "http://example.com/identity/",
					"region": "north",
					"links": {
						"self": "http://example.com/identity/v3/endpoints/6fedc0"
					},
					"service_id": "1b501a"
				},
				{
					"id": "6fedc0",
					"interface": "internal",
					"region": "south",
					"url": "http://example.com/identity/",
					"links": {
						"self": "http://example.com/identity/v3/endpoints/6fedc0"
					},
					"service_id": "1b501a"
				}
			],
			"links": {
				"self": "http://example.com/identity/v3/OS-EP-FILTER/projects/263fd9/endpoints",
				"previous": null,
				"next": null
			}
		}
		`)
	})

	count := 0
	err := projectendpoints.List(client.ServiceClient(), "project-id").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := projectendpoints.ExtractEndpoints(page)
		if err != nil {
			t.Errorf("Failed to extract endpoints: %v", err)
			return false, err
		}

		expected := []projectendpoints.Endpoint{
			{
				ID:           "6fedc0",
				Availability: gophercloud.AvailabilityPublic,
				Region:       "north",
				ServiceID:    "1b501a",
				URL:          "http://example.com/identity/",
			},
			{
				ID:           "6fedc0",
				Availability: gophercloud.AvailabilityInternal,
				Region:       "south",
				ServiceID:    "1b501a",
				URL:          "http://example.com/identity/",
			},
		}
		th.AssertDeepEquals(t, expected, actual)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestDeleteEndpoint(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/OS-EP-FILTER/projects/project-id/endpoints/endpoint-id", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	res := projectendpoints.Delete(context.TODO(), client.ServiceClient(), "project-id", "endpoint-id")
	th.AssertNoErr(t, res.Err)
}
