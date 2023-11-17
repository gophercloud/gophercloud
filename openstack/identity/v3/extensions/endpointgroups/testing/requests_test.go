package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/endpointgroups"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGetEndpointGroup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/OS-EP-FILTER/endpoint_groups/24", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		fmt.Fprintf(w, `
		{
			"endpoint_group": {
			    "id": "24",
			    "filters": {
				    "interface": "public",
					"service_id": "1234",
					"region_id": "5678"
			    },
			    "name": "endpointgroup1",
			    "description": "public endpoint group 1",
			    "links": {
					"self": "https://localhost:5000/v3/OS-EP-FILTER/endpoint_groups/24"
			    }
			}
		}
		`)
	})

	actual, err := endpointgroups.Get(client.ServiceClient(), "24").Extract()
	if err != nil {
		t.Fatalf("Failed to extract EndpointGroup: %v", err)
	}

	expected := &endpointgroups.EndpointGroup{
		ID: "24",
		Filters: map[string]interface{}{
			"interface":  "public",
			"service_id": "1234",
			"region_id":  "5678",
		},
		Name:        "endpointgroup1",
		Description: "public endpoint group 1",
	}
	th.AssertDeepEquals(t, expected, actual)
}

func TestListEndpointGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/OS-EP-FILTER/endpoint_groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprintf(w, `
		{
			"endpoint_groups": [
				{
				    "id": "24",
				    "filters": {
						"interface": "public",
						"service_id": "1234",
						"region_id": "5678",
				    },
				    "name": "endpointgroup1",
				    "description": "public endpoint group 1",
				    "links": {
						"self": "https://localhost:5000/v3/OS-EP-FILTER/endpoint_groups/24"
				    }
				},
				{
				    "id": "25",
				    "filters": {
						"interface": "internal"
				    },
				    "name": "endpointgroup2",
				    "description": "internal endpoint group 1",
				    "links": {
						"self": "https://localhost:5000/v3/OS-EP-FILTER/endpoint_groups/25"
				    }
				}
			]
		}
		`)
	})

	count := 0
	endpointgroups.List(client.ServiceClient(), endpointgroups.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := endpointgroups.ExtractEndpointGroups(page)
		if err != nil {
			t.Errorf("Failed to extract EndpointGroups: %v", err)
			return false, err
		}

		expected := []endpointgroups.EndpointGroup{
			{
				ID: "24",
				Filters: map[string]interface{}{
					"interface":  "public",
					"service_id": "1234",
					"region_id":  "5678",
				},
				Name:        "endpointgroup1",
				Description: "public endpoint group 1",
			},
			{
				ID: "25",
				Filters: map[string]interface{}{
					"interface": "internal",
				},
				Name:        "endpointgroup2",
				Description: "internal endpoint group 1",
			},
		}
		th.AssertDeepEquals(t, expected, actual)
		return true, nil
	})
}
