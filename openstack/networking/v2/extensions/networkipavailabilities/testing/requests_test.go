package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/networkipavailabilities"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/network-ip-availabilities", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, NetworkIPAvailabilityListResult)
	})

	count := 0

	err := networkipavailabilities.List(fake.ServiceClient(fakeServer), networkipavailabilities.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/network-ip-availabilities/cf11ab78-2302-49fa-870f-851a08c7afb8", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, NetworkIPAvailabilityGetResult)
	})

	s, err := networkipavailabilities.Get(context.TODO(), fake.ServiceClient(fakeServer), "cf11ab78-2302-49fa-870f-851a08c7afb8").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "cf11ab78-2302-49fa-870f-851a08c7afb8", s.NetworkID)
	th.AssertEquals(t, "public", s.NetworkName)
	th.AssertEquals(t, "424e7cf0243c468ca61732ba45973b3e", s.ProjectID)
	th.AssertEquals(t, "424e7cf0243c468ca61732ba45973b3e", s.TenantID)
	th.AssertEquals(t, "253", s.TotalIPs)
	th.AssertEquals(t, "3", s.UsedIPs)
	th.AssertDeepEquals(t, []networkipavailabilities.SubnetIPAvailability{
		{
			SubnetID:   "4afe6e5f-9649-40db-b18f-64c7ead942bd",
			SubnetName: "public-subnet",
			CIDR:       "203.0.113.0/24",
			IPVersion:  int(gophercloud.IPv4),
			TotalIPs:   "253",
			UsedIPs:    "3",
		},
	}, s.SubnetIPAvailabilities)
}
