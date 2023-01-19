package testing

import (
	"testing"

	"github.com/bizflycloud/gophercloud/openstack/orchestration/v1/buildinfo"
	th "github.com/bizflycloud/gophercloud/testhelper"
	fake "github.com/bizflycloud/gophercloud/testhelper/client"
)

func TestGetTemplate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t, GetOutput)

	actual, err := buildinfo.Get(fake.ServiceClient()).Extract()
	th.AssertNoErr(t, err)

	expected := GetExpected
	th.AssertDeepEquals(t, expected, actual)
}
