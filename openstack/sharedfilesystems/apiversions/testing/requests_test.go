package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/apiversions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAPIVersions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	allVersions, err := apiversions.List(client.ServiceClient()).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := apiversions.ExtractAPIVersions(allVersions)
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, ManilaAllAPIVersionResults, actual)
}

func TestGetAPIVersion(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	actual, err := apiversions.Get(context.TODO(), client.ServiceClient(), "v2").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, ManilaAPIVersion2Result, *actual)
}

func TestGetNoAPIVersion(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetNoResponse(t)

	_, err := apiversions.Get(context.TODO(), client.ServiceClient(), "v2").Extract()
	th.AssertEquals(t, err.Error(), "Unable to find requested API version")
}

func TestGetMultipleAPIVersion(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetMultipleResponses(t)

	_, err := apiversions.Get(context.TODO(), client.ServiceClient(), "v2").Extract()
	th.AssertEquals(t, err.Error(), "Found 2 API versions")
}
