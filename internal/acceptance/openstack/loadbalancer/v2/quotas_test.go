//go:build acceptance || networking || loadbalancer || quotas

package v2

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/quotas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestQuotasGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	quotasInfo, err := quotas.Get(context.TODO(), client, os.Getenv("OS_PROJECT_NAME")).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotasInfo)
}

func TestQuotasUpdate(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewLoadBalancerV2Client()
	th.AssertNoErr(t, err)

	originalQuotas, err := quotas.Get(context.TODO(), client, os.Getenv("OS_PROJECT_NAME")).Extract()
	th.AssertNoErr(t, err)

	var quotaUpdateOpts = quotas.UpdateOpts{
		Loadbalancer:  gophercloud.IntToPointer(25),
		Listener:      gophercloud.IntToPointer(45),
		Member:        gophercloud.IntToPointer(205),
		Pool:          gophercloud.IntToPointer(25),
		Healthmonitor: gophercloud.IntToPointer(5),
	}
	// L7 parameters are only supported in microversion v2.19 introduced in victoria
	if clients.IsCurrentAbove(t, "stable/ussuri") {
		quotaUpdateOpts.L7Policy = gophercloud.IntToPointer(55)
		quotaUpdateOpts.L7Rule = gophercloud.IntToPointer(105)
	}

	newQuotas, err := quotas.Update(context.TODO(), client, os.Getenv("OS_PROJECT_NAME"), quotaUpdateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newQuotas)

	if reflect.DeepEqual(originalQuotas, newQuotas) {
		log.Fatal("Original and New Loadbalancer Quotas are the same")
	}

	var restoredQuotaUpdate = quotas.UpdateOpts{
		Loadbalancer:  &originalQuotas.Loadbalancer,
		Listener:      &originalQuotas.Listener,
		Member:        &originalQuotas.Member,
		Pool:          &originalQuotas.Pool,
		Healthmonitor: &originalQuotas.Healthmonitor,
	}
	// L7 parameters are only supported in microversion v2.19 introduced in victoria
	if clients.IsCurrentAbove(t, "stable/ussuri") {
		restoredQuotaUpdate.L7Policy = &originalQuotas.L7Policy
		restoredQuotaUpdate.L7Rule = &originalQuotas.L7Rule
	}

	// Restore original quotas.
	restoredQuotas, err := quotas.Update(context.TODO(), client, os.Getenv("OS_PROJECT_NAME"), restoredQuotaUpdate).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, originalQuotas, restoredQuotas)

	tools.PrintResource(t, restoredQuotas)
}
