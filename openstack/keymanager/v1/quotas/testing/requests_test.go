package testing

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/quotas"
	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGet_1(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/quotas", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetResponseRaw_1)
	})

	q, err := quotas.Get(client.ServiceClient()).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, q, &GetResponse)
}

func TestListQuotas(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/project-quotas", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponseRaw_1)
	})

	count := 0
	err := quotas.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := quotas.ExtractQuotas(page)
		th.AssertNoErr(t, err)

		th.AssertDeepEquals(t, ExpectedQuotasSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, count, 1)
}

func TestListOrdersAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	th.Mux.HandleFunc("/project-quotas", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponseRaw_1)
	})

	allPages, err := quotas.List(client.ServiceClient(), nil).AllPages()
	th.AssertNoErr(t, err)
	actual, err := quotas.ExtractQuotas(allPages)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedQuotasSlice, actual)
}

func TestGetProjectQuota_1(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/project-quotas/0a73845280574ad389c292f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, GetProjectResponseRaw_1)
	})

	q, err := quotas.GetProjectQuota(client.ServiceClient(), "0a73845280574ad389c292f6a74afa76").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, q, &GetResponse)
}

func TestUpdate_1(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/project-quotas/0a73845280574ad389c292f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestJSONRequest(t, r, `
		{
			"project_quotas": {
				"secrets": 10,
				"orders": null,
				"containers": 14,
				"consumers": 15,
				"cas": null
			}
		}`)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

	})

	err := quotas.Update(client.ServiceClient(), "0a73845280574ad389c292f6a74afa76", quotas.UpdateOpts{
		Secrets:    gophercloud.IntToPointer(10),
		Orders:     nil,
		Containers: gophercloud.IntToPointer(14),
		Consumers:  gophercloud.IntToPointer(15),
		CAS:        nil,
	}).Err

	th.AssertNoErr(t, err)
}

func TestDelete_1(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/project-quotas/0a73845280574ad389c292f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)

	})

	err := quotas.Delete(client.ServiceClient(), "0a73845280574ad389c292f6a74afa76").Err

	th.AssertNoErr(t, err)
}
