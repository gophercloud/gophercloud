package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/clustering/v1/clusterpolicies"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListActions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListSuccessfully(t)

	pageCount := 0
	clusterID := "7d85f602-a948-4a30-afd4-e84f47471c15"
	err := clusterpolicies.List(fake.ServiceClient(), clusterID, clusterpolicies.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		pageCount++
		actual, err := clusterpolicies.ExtractClusterPolicies(page)
		th.AssertNoErr(t, err)

		th.AssertDeepEquals(t, ExpectedListClusterPolicy, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if pageCount != 1 {
		t.Errorf("Expected 1 page, got %d", pageCount)
	}
}
