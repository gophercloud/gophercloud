package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/addressscopes"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/address-scopes", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, AddressScopesListResult)
	})

	count := 0

	addressscopes.List(fake.ServiceClient(), addressscopes.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := addressscopes.ExtractAddressScopes(page)
		if err != nil {
			t.Errorf("Failed to extract addressscopes: %v", err)
			return false, nil
		}

		expected := []addressscopes.AddressScope{
			AddressScope1,
			AddressScope2,
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
