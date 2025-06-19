package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/limits"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetEnforcementModel(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetEnforcementModelSuccessfully(t, fakeServer)

	actual, err := limits.GetEnforcementModel(context.TODO(), client.ServiceClient(fakeServer)).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Model, *actual)
}

func TestListLimits(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListLimitsSuccessfully(t, fakeServer)

	count := 0
	err := limits.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListLimitsSuccessfully(t, fakeServer)

	allPages, err := limits.List(client.ServiceClient(fakeServer), nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := limits.ExtractLimits(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedLimitsSlice, actual)
}

func TestCreateLimits(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateLimitSuccessfully(t, fakeServer)

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

	actual, err := limits.BatchCreate(context.TODO(), client.ServiceClient(fakeServer), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedLimitsSlice, actual)
}

func TestGetLimit(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetLimitSuccessfully(t, fakeServer)

	actual, err := limits.Get(context.TODO(), client.ServiceClient(fakeServer), "25a04c7a065c430590881c646cdcdd58").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, FirstLimit, *actual)
}

func TestUpdateLimit(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateLimitSuccessfully(t, fakeServer)

	var description = "Number of snapshots for project 3a705b9f56bb439381b43c4fe59dccce"
	var resourceLimit = 5
	updateOpts := limits.UpdateOpts{
		Description:   &description,
		ResourceLimit: &resourceLimit,
	}

	actual, err := limits.Update(context.TODO(), client.ServiceClient(fakeServer), "3229b3849f584faea483d6851f7aab05", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondLimitUpdated, *actual)
}

func TestDeleteLimit(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteLimitSuccessfully(t, fakeServer)

	err := limits.Delete(context.TODO(), client.ServiceClient(fakeServer), "3229b3849f584faea483d6851f7aab05").ExtractErr()
	th.AssertNoErr(t, err)
}
