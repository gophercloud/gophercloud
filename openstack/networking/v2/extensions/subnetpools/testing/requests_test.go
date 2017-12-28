package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/subnetpools"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/subnetpools", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, SubnetPoolsListResult)
	})

	count := 0

	subnetpools.List(fake.ServiceClient(), subnetpools.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := subnetpools.ExtractSubnetPools(page)
		if err != nil {
			t.Errorf("Failed to extract subnetpools: %v", err)
			return false, nil
		}

		expected := []subnetpools.SubnetPool{
			SubnetPool1,
			SubnetPool2,
			SubnetPool3,
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
