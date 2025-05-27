package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/addressgroups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/address-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, AddressGroupListResponse)
	})

	count := 0

	err := addressgroups.List(fake.ServiceClient(fakeServer), addressgroups.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := addressgroups.ExtractGroups(page)
		if err != nil {
			t.Errorf("Failed to extract address groups: %v", err)
			return false, err
		}

		expected := []addressgroups.AddressGroup{
			{
				Description: "",
				ID:          "8722e0e0-9cc9-4490-9660-8c9a5732fbb0",
				Name:        "ADDR_GP_1",
				ProjectID:   "45977fa2dbd7482098dd68d0d8970117",
				Addresses: []string{
					"132.168.4.12/24",
				},
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/address-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, AddressGroupCreateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, AddressGroupCreateResponse)
	})

	opts := addressgroups.CreateOpts{
		Name: "ADDR_GP_1",
		Addresses: []string{
			"132.168.4.12/24",
		},
	}
	_, err := addressgroups.Create(context.TODO(), fake.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)
}

func TestRequiredCreateOpts(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	_, err := addressgroups.Create(context.TODO(), fake.ServiceClient(fakeServer), addressgroups.CreateOpts{Name: "ADDR_GP_1"}).Extract()
	if err == nil {
		t.Fatalf("Expected error, got none")
	}
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/address-groups/8722e0e0-9cc9-4490-9660-8c9a5732fbb0", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, AddressGroupGetResponse)
	})

	sr, err := addressgroups.Get(context.TODO(), fake.ServiceClient(fakeServer), "8722e0e0-9cc9-4490-9660-8c9a5732fbb0").Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "", sr.Description)
	th.AssertEquals(t, "8722e0e0-9cc9-4490-9660-8c9a5732fbb0", sr.ID)
	th.AssertEquals(t, "45977fa2dbd7482098dd68d0d8970117", sr.ProjectID)
	th.CheckDeepEquals(t, []string{"132.168.4.12/24"}, sr.Addresses)
	th.AssertEquals(t, "ADDR_GP_1", sr.Name)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/address-groups/8722e0e0-9cc9-4490-9660-8c9a5732fbb0",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, AddressGroupUpdateRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, AddressGroupUpdateResponse)
		})

	name := "ADDR_GP_2"
	description := "new description"
	opts := addressgroups.UpdateOpts{
		Name:        &name,
		Description: &description,
	}
	ag, err := addressgroups.Update(context.TODO(), fake.ServiceClient(fakeServer), "8722e0e0-9cc9-4490-9660-8c9a5732fbb0", opts).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, []string{"192.168.4.1/32"}, ag.Addresses)
	th.AssertEquals(t, "new description", ag.Description)
	th.AssertEquals(t, "8722e0e0-9cc9-4490-9660-8c9a5732fbb0", ag.ID)
	th.AssertEquals(t, "45977fa2dbd7482098dd68d0d8970117", ag.ProjectID)
	th.AssertEquals(t, "ADDR_GP_2", ag.Name)
}

func TestAddAddresses(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/address-groups/8722e0e0-9cc9-4490-9660-8c9a5732fbb0/add_addresses",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, AddressGroupAddAddressesRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, AddressGroupAddAddressesResponse)
		})

	opts := addressgroups.UpdateAddressesOpts{
		Addresses: []string{"192.168.4.1/32"},
	}
	ag, err := addressgroups.AddAddresses(context.TODO(), fake.ServiceClient(fakeServer), "8722e0e0-9cc9-4490-9660-8c9a5732fbb0", opts).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, []string{"132.168.4.12/24", "192.168.4.1/32"}, ag.Addresses)
	th.AssertEquals(t, "original description", ag.Description)
	th.AssertEquals(t, "8722e0e0-9cc9-4490-9660-8c9a5732fbb0", ag.ID)
}

func TestRemoveAddresses(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/address-groups/8722e0e0-9cc9-4490-9660-8c9a5732fbb0/remove_addresses",
		func(w http.ResponseWriter, r *http.Request) {
			th.TestMethod(t, r, "PUT")
			th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
			th.TestHeader(t, r, "Content-Type", "application/json")
			th.TestHeader(t, r, "Accept", "application/json")
			th.TestJSONRequest(t, r, AddressGroupRemoveAddressesRequest)

			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			fmt.Fprintf(w, AddressGroupRemoveAddressesResponse)
		})

	opts := addressgroups.UpdateAddressesOpts{
		Addresses: []string{"192.168.4.1/32"},
	}
	ag, err := addressgroups.RemoveAddresses(context.TODO(), fake.ServiceClient(fakeServer), "8722e0e0-9cc9-4490-9660-8c9a5732fbb0", opts).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, []string{"132.168.4.12/24"}, ag.Addresses)
	th.AssertEquals(t, "original description", ag.Description)
	th.AssertEquals(t, "8722e0e0-9cc9-4490-9660-8c9a5732fbb0", ag.ID)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	fakeServer.Mux.HandleFunc("/v2.0/address-groups/4ec89087-d057-4e2c-911f-60a3b47ee304", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.WriteHeader(http.StatusNoContent)
	})

	err := addressgroups.Delete(context.TODO(), fake.ServiceClient(fakeServer), "4ec89087-d057-4e2c-911f-60a3b47ee304").ExtractErr()
	th.AssertNoErr(t, err)
}
