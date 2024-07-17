package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/apiversions"
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

	th.AssertDeepEquals(t, NovaAllAPIVersionResults, actual)
}

func TestGetAPIVersion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetResponse(t, fakeServer)

	actual, err := apiversions.Get(context.TODO(), client.ServiceClient(fakeServer), "v2.1").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, NovaAPIVersion21Result, *actual)
}

func TestGetMultipleAPIVersion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockGetMultipleResponses(t, fakeServer)

	_, err := apiversions.Get(context.TODO(), client.ServiceClient(fakeServer), "v3").Extract()
	th.AssertEquals(t, err.Error(), "Unable to find requested API version")
}
