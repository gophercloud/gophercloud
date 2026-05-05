package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/usages"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetUsagesSuccess(t, fakeServer)

	actual, err := usages.Get(context.TODO(), client.ServiceClient(fakeServer), usages.GetOpts{
		ProjectID: ProjectID,
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedUsages, *actual)
}

func TestGetWithUserSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetUsagesWithUserSuccess(t, fakeServer)

	actual, err := usages.Get(context.TODO(), client.ServiceClient(fakeServer), usages.GetOpts{
		ProjectID: ProjectID,
		UserID:    UserID,
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedUsages, *actual)
}

func TestGetEmptySuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetEmptyUsagesSuccess(t, fakeServer)

	actual, err := usages.Get(context.TODO(), client.ServiceClient(fakeServer), usages.GetOpts{
		ProjectID: ProjectID,
	}).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedEmptyUsages, *actual)
}

func TestGetPre138Success(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetUsagesPre138Success(t, fakeServer)

	actual, err := usages.Get(context.TODO(), client.ServiceClient(fakeServer), usages.GetOpts{
		ProjectID: ProjectID,
	}).ExtractPre138()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedUsagesPre138, *actual)
}

func TestGetPre138WithUserSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetUsagesPre138WithUserSuccess(t, fakeServer)

	actual, err := usages.Get(context.TODO(), client.ServiceClient(fakeServer), usages.GetOpts{
		ProjectID: ProjectID,
		UserID:    UserID,
	}).ExtractPre138()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedUsagesPre138, *actual)
}

func TestGetPre138EmptySuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetEmptyUsagesPre138Success(t, fakeServer)

	actual, err := usages.Get(context.TODO(), client.ServiceClient(fakeServer), usages.GetOpts{
		ProjectID: ProjectID,
	}).ExtractPre138()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedEmptyUsagesPre138, *actual)
}
