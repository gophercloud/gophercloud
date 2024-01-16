package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/limits"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestGetEnforcementModel(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetEnforcementModelSuccessfully(t)

	actual, err := limits.GetEnforcementModel(client.ServiceClient()).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Model, *actual)
}

func TestListLimits(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListLimitsSuccessfully(t)

	count := 0
	err := limits.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := limits.ExtractLimits(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedLimitsSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListLimitsAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListLimitsSuccessfully(t)

	allPages, err := limits.List(client.ServiceClient(), nil).AllPages()
	th.AssertNoErr(t, err)
	actual, err := limits.ExtractLimits(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedLimitsSlice, actual)
}

func TestCreateLimits(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateLimitSuccessfully(t)

	createOpts := limits.BatchCreateOpts{
		limits.CreateOpts{
			ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
			ProjectID:     "3a705b9f56bb439381b43c4fe59dccce",
			RegionID:      "RegionOne",
			ResourceName:  "snapshot",
			ResourceLimit: 5,
		},
		limits.CreateOpts{
			ServiceID:     "9408080f1970482aa0e38bc2d4ea34b7",
			DomainID:      "edbafc92be354ffa977c58aa79c7bdb2",
			ResourceName:  "volume",
			ResourceLimit: 11,
			Description:   "Number of volumes for project 3a705b9f56bb439381b43c4fe59dccce",
		},
	}

	actual, err := limits.BatchCreate(client.ServiceClient(), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedLimitsSlice, actual)
}

func TestGetLimit(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetLimitSuccessfully(t)

	actual, err := limits.Get(client.ServiceClient(), "25a04c7a065c430590881c646cdcdd58").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, FirstLimit, *actual)
}

func TestUpdateLimit(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateLimitSuccessfully(t)

	description := "Number of snapshots for project 3a705b9f56bb439381b43c4fe59dccce"
	resourceLimit := 5
	updateOpts := limits.UpdateOpts{
		Description:   &description,
		ResourceLimit: &resourceLimit,
	}

	actual, err := limits.Update(client.ServiceClient(), "3229b3849f584faea483d6851f7aab05", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondLimitUpdated, *actual)
}

func TestDeleteLimit(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteLimitSuccessfully(t)

	err := limits.Delete(client.ServiceClient(), "3229b3849f584faea483d6851f7aab05").ExtractErr()
	th.AssertNoErr(t, err)
}
