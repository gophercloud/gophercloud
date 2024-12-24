package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/quotas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGet(t)

	expectedQuotas := []quotas.QuotaDetail{
		{Resource: "instances", Limit: 15, InUse: 5, Reserved: 0},
		{Resource: "backups", Limit: 50, InUse: 2, Reserved: 0},
		{Resource: "volumes", Limit: 40, InUse: 1, Reserved: 0},
	}

	actual, err := quotas.Get(context.TODO(), client.ServiceClient(), "e131f89a-c1d8-11ef-bfaa-370c246e2439").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, expectedQuotas, actual)
}
