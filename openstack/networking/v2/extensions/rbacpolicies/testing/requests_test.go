package testing

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	fake "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/common"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/rbacpolicies"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
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

		fmt.Fprint(w, CreateResponse)
	})

	options := rbacpolicies.CreateOpts{
		Action:       rbacpolicies.ActionAccessShared,
		ObjectType:   "network",
		TargetTenant: "6e547a3bcfe44702889fdeff3c3520c3",
		ObjectID:     "240d22bf-bd17-4238-9758-25f72610ecdc",
	}
	rbacResult, err := rbacpolicies.Create(context.TODO(), fake.ServiceClient(), options).Extract()
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

		fmt.Fprint(w, GetResponse)
	})

	n, err := rbacpolicies.Get(context.TODO(), fake.ServiceClient(), "2cf7523a-93b5-4e69-9360-6c6bf986bb7c").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &rbacPolicy1, n)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/rbac-policies", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})

	client := fake.ServiceClient()
	count := 0

	err := rbacpolicies.List(client, rbacpolicies.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := rbacpolicies.ExtractRBACPolicies(page)
		if err != nil {
			t.Errorf("Failed to extract rbac policies: %v", err)
			return false, err
		}

		th.CheckDeepEquals(t, ExpectedRBACPoliciesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestListWithAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/rbac-policies", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, ListResponse)
	})

	client := fake.ServiceClient()

	type newRBACPolicy struct {
		rbacpolicies.RBACPolicy
	}

	var allRBACpolicies []newRBACPolicy

	allPages, err := rbacpolicies.List(client, rbacpolicies.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	err = rbacpolicies.ExtractRBACPolicesInto(allPages, &allRBACpolicies)
	th.AssertNoErr(t, err)

	th.AssertEquals(t, allRBACpolicies[0].ObjectType, "network")
	th.AssertEquals(t, allRBACpolicies[0].Action, rbacpolicies.ActionAccessShared)

	th.AssertEquals(t, allRBACpolicies[1].ProjectID, "1ae27ce0a2a54cc6ae06dc62dd0ec832")
	th.AssertEquals(t, allRBACpolicies[1].TargetTenant, "1a547a3bcfe44702889fdeff3c3520c3")

}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/rbac-policies/71d55b18-d2f8-4c76-a5e6-e0a3dd114361", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "DELETE")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		w.WriteHeader(http.StatusNoContent)
	})

	res := rbacpolicies.Delete(context.TODO(), fake.ServiceClient(), "71d55b18-d2f8-4c76-a5e6-e0a3dd114361").ExtractErr()
	th.AssertNoErr(t, res)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/rbac-policies/2cf7523a-93b5-4e69-9360-6c6bf986bb7c", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PUT")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		th.TestHeader(t, r, "Content-Type", "application/json")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestJSONRequest(t, r, UpdateRequest)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprint(w, UpdateResponse)
	})

	options := rbacpolicies.UpdateOpts{TargetTenant: "9d766060b6354c9e8e2da44cab0e8f38"}
	rbacResult, err := rbacpolicies.Update(context.TODO(), fake.ServiceClient(), "2cf7523a-93b5-4e69-9360-6c6bf986bb7c", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, rbacResult.TargetTenant, "9d766060b6354c9e8e2da44cab0e8f38")
	th.AssertEquals(t, rbacResult.ID, "2cf7523a-93b5-4e69-9360-6c6bf986bb7c")
}
