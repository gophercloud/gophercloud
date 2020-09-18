// +build acceptance networking loadbalancer quotas

package v2

import (
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/quotas"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestQuotasGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	quotasInfo, err := quotas.Get(client, os.Getenv("OS_PROJECT_NAME")).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotasInfo)
}

func TestQuotasUpdate(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	originalQuotas, err := quotas.Get(client, os.Getenv("OS_PROJECT_NAME")).Extract()
	th.AssertNoErr(t, err)

	newQuotas, err := quotas.Update(client, os.Getenv("OS_PROJECT_NAME"), quotaUpdateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newQuotas)

	if reflect.DeepEqual(originalQuotas, newQuotas) {
		log.Fatal("Original and New Loadbalancer Quotas are the same")
	}

	// Restore original quotas.
	restoredQuotas, err := quotas.Update(client, os.Getenv("OS_PROJECT_NAME"), quotas.UpdateOpts{
		Loadbalancer:  &originalQuotas.Loadbalancer,
		Listener:      &originalQuotas.Listener,
		Member:        &originalQuotas.Member,
		Pool:          &originalQuotas.Pool,
		Healthmonitor: &originalQuotas.Healthmonitor,
		L7Policy:      &originalQuotas.L7Policy,
		L7Rule:        &originalQuotas.L7Rule,
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, originalQuotas, restoredQuotas)

	tools.PrintResource(t, restoredQuotas)
}
