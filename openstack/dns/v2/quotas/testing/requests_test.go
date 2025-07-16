package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/quotas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer)

	actual, err := quotas.Get(context.TODO(), client.ServiceClient(fakeServer), "a86dba58-0043-4cc6-a1bb-69d5e86f3ca3").Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Quota, actual)
}

func TestUpdate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleUpdateSuccessfully(t, fakeServer)

	zones := 100
	updateOpts := quotas.UpdateOpts{
		Zones: &zones,
	}

	actual, err := quotas.Update(context.TODO(), client.ServiceClient(fakeServer), "a86dba58-0043-4cc6-a1bb-69d5e86f3ca3", updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, Quota, actual)
}
