package testing

import (
	"context"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/apiversions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAPIVersions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockListResponse(t, fakeServer)

	allVersions, err := apiversions.List(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := apiversions.ExtractAPIVersions(allVersions)
	th.AssertNoErr(t, err)
	fmt.Println(actual)
	th.AssertDeepEquals(t, MagnumAllAPIVersionResults, actual)
}

func TestGetAPIVersion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetResponse(t, fakeServer)

	actual, err := apiversions.Get(context.TODO(), client.ServiceClient(fakeServer), "v1").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, MagnumAPIVersion1Result, *actual)
}
