package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/networkipavailabilities"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/network-ip-availabilities", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, NetworkIPAvailabilityListResult)
	})

	count := 0

	networkipavailabilities.List(fake.ServiceClient(), networkipavailabilities.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := networkipavailabilities.ExtractNetworkIPAvailabilities(page)
		if err != nil {
			t.Errorf("Failed to extract network IP availabilities: %v", err)
			return false, nil
		}

		expected := []networkipavailabilities.NetworkIPAvailability{
			NetworkIPAvailability1,
			NetworkIPAvailability2,
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}
