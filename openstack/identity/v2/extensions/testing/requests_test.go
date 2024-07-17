package testing

import (
	"context"
	"testing"

	common "github.com/gophercloud/gophercloud/v2/openstack/common/extensions/testing"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/extensions"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListExtensionsSuccessfully(t, fakeServer)

	count := 0
	err := extensions.List(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := extensions.ExtractExtensions(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, common.ExpectedExtensions, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	common.HandleGetExtensionSuccessfully(t, fakeServer)

	actual, err := extensions.Get(context.TODO(), client.ServiceClient(fakeServer), "agent").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, common.SingleExtension, actual)
}
