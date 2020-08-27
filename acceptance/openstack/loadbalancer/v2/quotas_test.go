package v2

import (
	"os"
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
