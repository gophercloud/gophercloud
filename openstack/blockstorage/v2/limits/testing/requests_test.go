package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/limits"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	actual, err := limits.Get(context.TODO(), client.ServiceClient()).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &LimitsResult, actual)
}
