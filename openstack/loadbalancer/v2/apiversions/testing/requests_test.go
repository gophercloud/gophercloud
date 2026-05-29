package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/apiversions"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListVersions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockListResponse(t, fakeServer)

	allVersions, err := apiversions.List(client.ServiceClient(fakeServer)).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := apiversions.ExtractAPIVersions(allVersions)
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, OctaviaAllAPIVersionResults, actual)
}
