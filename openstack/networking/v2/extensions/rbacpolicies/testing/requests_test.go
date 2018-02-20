package testing

import (
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/rbacpolicies"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/rbac-policies", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, CreateRequest)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		fmt.Fprintf(w, CreateResponse)
	})

	options := rbacpolicies.CreateOpts{
		Action:       rbacpolicies.ActionAccessShared,
		ObjectType:   "network",
		TargetTenant: "6e547a3bcfe44702889fdeff3c3520c3",
		ObjectID:     "240d22bf-bd17-4238-9758-25f72610ecdc",
	}
	rbacResult, err := rbacpolicies.Create(fake.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, &rbacPolicy1, rbacResult)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/rbac-policies/2cf7523a-93b5-4e69-9360-6c6bf986bb7c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintf(w, GetResponse)
	})

	n, err := rbacpolicies.Get(fake.ServiceClient(), "2cf7523a-93b5-4e69-9360-6c6bf986bb7c").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &rbacPolicy1, n)
}
