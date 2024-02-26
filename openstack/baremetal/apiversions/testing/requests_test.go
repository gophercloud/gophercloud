package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/apiversions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAPIVersions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	actual, err := apiversions.List(context.TODO(), client.ServiceClient()).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, IronicAllAPIVersionResults, *actual)
}

func TestGetAPIVersion(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	actual, err := apiversions.Get(context.TODO(), client.ServiceClient(), "v1").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, IronicAPIVersion1Result, *actual)
}
