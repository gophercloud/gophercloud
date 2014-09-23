// +build acceptance networking quotas

package extensions

import (
	"testing"

	base "github.com/rackspace/gophercloud/acceptance/openstack/networking/v2"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/quotas"
)

func TestQuotas(t *testing.T) {
	base.Setup()
	defer Teardown()

	setQuotas(t)
	getQuotas(t)
	resetQuotas(t)
}

func getQuotas(t *testing.T) {
	qs, err := quotas.Get(base.Client).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, qs.Subnet, 10)
	th.AssertEquals(t, qs.Router, 10)
	th.AssertEquals(t, qs.Port, 50)
	th.AssertEquals(t, qs.Network, 10)
	th.AssertEquals(t, qs.FloatingIP, 50)
}

func setQuotas(t *testing.T) {
	i10, i20, i30 := 10, 20, 30
	opts := UpdateOpts{Member: &i30, Pool: &i20, SecGroup: &i20, HealthMonitor: &i10}
	qs, err := quotas.Update(base.Client).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, qs.Member, 30)
	th.AssertEquals(t, qs.Pool, 20)
	th.AssertEquals(t, qs.SecGroup, 20)
	th.AssertEquals(t, qs.HealthMonitor, 10)
}

func resetQuotas(t *testing.T) {
	res := quotas.Reset(base.Client)
	th.AssertNoErr(t, res.Err)
}
