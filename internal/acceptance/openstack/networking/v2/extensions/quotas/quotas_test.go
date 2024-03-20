//go:build acceptance || networking || quotas

package quotas

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/quotas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestQuotasGet(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	quotasInfo, err := quotas.Get(context.TODO(), client, os.Getenv("OS_PROJECT_NAME")).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotasInfo)
}

func TestQuotasUpdate(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewNetworkV2Client()
	th.AssertNoErr(t, err)

	originalQuotas, err := quotas.Get(context.TODO(), client, os.Getenv("OS_PROJECT_NAME")).Extract()
	th.AssertNoErr(t, err)

	newQuotas, err := quotas.Update(context.TODO(), client, os.Getenv("OS_PROJECT_NAME"), updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newQuotas)

	if reflect.DeepEqual(originalQuotas, newQuotas) {
		log.Fatal("Original and New Networking Quotas are the same")
	}

	// Restore original quotas.
	restoredQuotas, err := quotas.Update(context.TODO(), client, os.Getenv("OS_PROJECT_NAME"), quotas.UpdateOpts{
		FloatingIP:        &originalQuotas.FloatingIP,
		Network:           &originalQuotas.Network,
		Port:              &originalQuotas.Port,
		RBACPolicy:        &originalQuotas.RBACPolicy,
		Router:            &originalQuotas.Router,
		SecurityGroup:     &originalQuotas.SecurityGroup,
		SecurityGroupRule: &originalQuotas.SecurityGroupRule,
		Subnet:            &originalQuotas.Subnet,
		SubnetPool:        &originalQuotas.SubnetPool,
	}).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, originalQuotas, restoredQuotas)

	tools.PrintResource(t, restoredQuotas)
}
