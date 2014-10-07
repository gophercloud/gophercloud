package tenants

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
)

const tokenID = "1234123412341234"

func TestListTenants(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/tenants", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", tokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `
{
  "tenants": [
    {
      "id": "1234",
      "name": "Red Team",
      "description": "The team that is red",
      "enabled": true
    },
    {
      "id": "9876",
      "name": "Blue Team",
      "description": "The team that is blue",
      "enabled": false
    }
  ]
}
    `)
	})

	client := &gophercloud.ServiceClient{
		Provider: &gophercloud.ProviderClient{TokenID: tokenID},
		Endpoint: th.Endpoint(),
	}

	count := 0
	err := List(client, nil).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := ExtractTenants(page)
		th.AssertNoErr(t, err)

		expected := []Tenant{
			Tenant{
				ID:          "1234",
				Name:        "Red Team",
				Description: "The team that is red",
				Enabled:     true,
			},
			Tenant{
				ID:          "9876",
				Name:        "Blue Team",
				Description: "The team that is blue",
				Enabled:     false,
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}
