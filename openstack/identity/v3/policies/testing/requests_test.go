package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/policies"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListPolicies(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListPoliciesSuccessfully(t)

	count := 0
	err := policies.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := policies.ExtractPolicies(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedPoliciesSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListPoliciesAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListPoliciesSuccessfully(t)

	allPages, err := policies.List(client.ServiceClient(), nil).AllPages()
	th.AssertNoErr(t, err)
	actual, err := policies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedPoliciesSlice, actual)
}

func TestListPoliciesWithFilter(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListPoliciesSuccessfully(t)

	listOpts := policies.ListOpts{
		Type: "application/json",
	}
	allPages, err := policies.List(client.ServiceClient(), listOpts).AllPages()
	th.AssertNoErr(t, err)
	actual, err := policies.ExtractPolicies(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, []policies.Policy{SecondPolicy}, actual)
}
