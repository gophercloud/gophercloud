package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	count := 0
	err := servergroups.List(client.ServiceClient(), &servergroups.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := servergroups.ExtractServerGroups(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedServerGroupList, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)

	actual, err := servergroups.Create(context.TODO(), client.ServiceClient(), servergroups.CreateOpts{
		Name:     "test",
		Policies: []string{"anti-affinity"},
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &ExpectedServerGroupCreate, actual)
}

func TestCreateMicroversion(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateMicroversionSuccessfully(t)

	result := servergroups.Create(context.TODO(), client.ServiceClient(), servergroups.CreateOpts{
		Name:   "test",
		Policy: policy,
		Rules:  ExpectedServerGroupCreateMicroversion.Rules,
	})

	actual, err := result.Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &ExpectedServerGroupCreateMicroversion, actual)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	actual, err := servergroups.Get(context.TODO(), client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &ExpectedServerGroupGet, actual)
}

func TestGetMicroversion(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetMicroversionSuccessfully(t)

	actual, err := servergroups.Get(context.TODO(), client.ServiceClient(), "4d8c3732-a248-40ed-bebc-539a6ffd25c0").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &ExpectedServerGroupGetMicroversion, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)

	err := servergroups.Delete(context.TODO(), client.ServiceClient(), "616fb98f-46ca-475e-917e-2563e5a8cd19").ExtractErr()
	th.AssertNoErr(t, err)
}
