package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/quotas/0a73845280574ad389c292f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, GetResponseRaw)
	})

	q, err := quotas.Get(context.TODO(), fake.ServiceClient(), "0a73845280574ad389c292f6a74afa76").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, q, &GetResponse)
}

func TestGetDetail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/quotas/0a73845280574ad389c292f6a74afa76/details.json", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, GetDetailedResponseRaw)
	})

	q, err := quotas.GetDetail(context.TODO(), fake.ServiceClient(), "0a73845280574ad389c292f6a74afa76").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, q, &GetDetailResponse)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/quotas/0a73845280574ad389c292f6a74afa76", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, UpdateRequestResponseRaw)
	})

	q, err := quotas.Update(context.TODO(), fake.ServiceClient(), "0a73845280574ad389c292f6a74afa76", quotas.UpdateOpts{
		FloatingIP:        gophercloud.IntToPointer(0),
		Network:           gophercloud.IntToPointer(-1),
		Port:              gophercloud.IntToPointer(5),
		RBACPolicy:        gophercloud.IntToPointer(10),
		Router:            gophercloud.IntToPointer(15),
		SecurityGroup:     gophercloud.IntToPointer(20),
		SecurityGroupRule: gophercloud.IntToPointer(-1),
		Subnet:            gophercloud.IntToPointer(25),
		SubnetPool:        gophercloud.IntToPointer(0),
		Trunk:             gophercloud.IntToPointer(5),
	}).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, q, &UpdateResponse)
}
