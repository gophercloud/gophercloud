package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/subnetpools"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/subnetpools", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, SubnetPoolsListResult)
	})

	count := 0

	err := subnetpools.List(fake.ServiceClient(fakeServer), subnetpools.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/subnetpools/0a738452-8057-4ad3-89c2-92f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, SubnetPoolGetResult)
	})

	s, err := subnetpools.Get(context.TODO(), fake.ServiceClient(fakeServer), "0a738452-8057-4ad3-89c2-92f6a74afa76").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "0a738452-8057-4ad3-89c2-92f6a74afa76", s.ID)
	th.AssertEquals(t, "my-ipv6-pool", s.Name)
	th.AssertEquals(t, 2, s.DefaultQuota)
	th.AssertEquals(t, "1e2b9857295a4a3e841809ef492812c5", s.TenantID)
	th.AssertEquals(t, "1e2b9857295a4a3e841809ef492812c5", s.ProjectID)
	th.AssertEquals(t, s.CreatedAt, time.Date(2018, 1, 1, 0, 0, 1, 0, time.UTC))
	th.AssertEquals(t, s.UpdatedAt, time.Date(2018, 1, 1, 0, 10, 10, 0, time.UTC))
	th.AssertDeepEquals(t, []string{
		"2001:db8::a3/64",
	}, s.Prefixes)
	th.AssertEquals(t, 64, s.DefaultPrefixLen)
	th.AssertEquals(t, 64, s.MinPrefixLen)
	th.AssertEquals(t, 128, s.MaxPrefixLen)
	th.AssertEquals(t, "", s.AddressScopeID)
	th.AssertEquals(t, 6, s.IPversion)
	th.AssertFalse(t, s.Shared)
	th.AssertEquals(t, "ipv6 prefixes", s.Description)
	th.AssertTrue(t, s.IsDefault)
	th.AssertEquals(t, 2, s.RevisionNumber)
}
func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/subnetpools", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, SubnetPoolCreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprint(w, SubnetPoolCreateResult)
	})

	opts := subnetpools.CreateOpts{
		Name: "my_ipv4_pool",
		Prefixes: []string{
			"10.10.0.0/16",
			"10.11.11.0/24",
		},
		MinPrefixLen:   25,
		MaxPrefixLen:   30,
		AddressScopeID: "3d4e2e2a-552b-42ad-a16d-820bbf3edaf3",
		Description:    "ipv4 prefixes",
	}
	s, err := subnetpools.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "my_ipv4_pool", s.Name)
	th.AssertDeepEquals(t, []string{
		"10.10.0.0/16",
		"10.11.11.0/24",
	}, s.Prefixes)
	th.AssertEquals(t, 25, s.MinPrefixLen)
	th.AssertEquals(t, 30, s.MaxPrefixLen)
	th.AssertEquals(t, "3d4e2e2a-552b-42ad-a16d-820bbf3edaf3", s.AddressScopeID)
	th.AssertEquals(t, "ipv4 prefixes", s.Description)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/subnetpools/099546ca-788d-41e5-a76d-17d8cd282d3e", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, SubnetPoolUpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, SubnetPoolUpdateResponse)
	})

	nullString := ""
	nullInt := 0
	updateOpts := subnetpools.UpdateOpts{
		Name: "new_subnetpool_name",
		Prefixes: []string{
			"10.11.12.0/24",
			"10.24.0.0/16",
		},
		MaxPrefixLen:   16,
		AddressScopeID: &nullString,
		DefaultQuota:   &nullInt,
		Description:    &nullString,
	}
	n, err := subnetpools.Update(context.TODO(), fake.ServiceClient(fakeServer), "099546ca-788d-41e5-a76d-17d8cd282d3e", updateOpts).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "new_subnetpool_name", n.Name)
	th.AssertDeepEquals(t, []string{
		"10.8.0.0/16",
		"10.11.12.0/24",
		"10.24.0.0/16",
	}, n.Prefixes)
	th.AssertEquals(t, 16, n.MaxPrefixLen)
	th.AssertEquals(t, "099546ca-788d-41e5-a76d-17d8cd282d3e", n.ID)
	th.AssertEquals(t, "", n.AddressScopeID)
	th.AssertEquals(t, 0, n.DefaultQuota)
	th.AssertEquals(t, "", n.Description)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/subnetpools/099546ca-788d-41e5-a76d-17d8cd282d3e", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := subnetpools.Delete(context.TODO(), fake.ServiceClient(fakeServer), "099546ca-788d-41e5-a76d-17d8cd282d3e")
	th.AssertNoErr(t, res.Err)
}
