package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/orchestration/v1/buildinfo"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetTemplate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer, GetOutput)

	actual, err := buildinfo.Get(context.TODO(), client.ServiceClient(fakeServer)).Extract()
	th.AssertNoErr(t, err)

	expected := GetExpected
	th.AssertDeepEquals(t, expected, actual)
}
