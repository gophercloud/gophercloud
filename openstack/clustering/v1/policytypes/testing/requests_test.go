package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/clustering/v1/policytypes"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListPolicyTypes(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandlePolicyTypeList(t)

	count := 0
	err := policytypes.List(fake.ServiceClient()).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := policytypes.ExtractPolicyTypes(page)
		if err != nil {
			t.Errorf("Failed to extract policy types: %v", err)
			return false, err
		}
		th.AssertDeepEquals(t, ExpectedPolicyTypes, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGetPolicyType(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandlePolicyTypeGet(t)

	actual, err := policytypes.Get(context.TODO(), fake.ServiceClient(), FakePolicyTypetoGet).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, ExpectedPolicyTypeDetail, actual)
}
