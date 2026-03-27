package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/traits"

	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListTraitsAll(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListTraitsAll(t, fakeServer)

	count := 0
	err := traits.List(client.ServiceClient(fakeServer), traits.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := traits.ExtractTraits(page)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, ExpectedTraitsListResultAll, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	th.AssertEquals(t, 1, count)
}

func TestListTraitsFilteredName(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListTraitsFilteredName(t, fakeServer)

	count := 0
	err := traits.List(client.ServiceClient(fakeServer), traits.ListOpts{Name: "startswith:CUSTOM"}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := traits.ExtractTraits(page)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, ExpectedTraitsListFilteredNameResult, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	th.AssertEquals(t, 1, count)
}

func TestListTraitsFilteredAssociated(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListTraitsFilteredAssociated(t, fakeServer)

	count := 0
	associated := true
	err := traits.List(client.ServiceClient(fakeServer), traits.ListOpts{Associated: &associated}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++

		actual, err := traits.ExtractTraits(page)
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, ExpectedTraitsListFilteredAssociatedResult, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	th.AssertEquals(t, 1, count)
}

func TestGetTraitSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetTraitSuccess(t, fakeServer)

	err := traits.Get(context.TODO(), client.ServiceClient(fakeServer), PresentTrait).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetTraitNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetTraitNotFound(t, fakeServer)

	err := traits.Get(context.TODO(), client.ServiceClient(fakeServer), AbsentTrait).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}
