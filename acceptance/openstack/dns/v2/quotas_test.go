//go:build acceptance
// +build acceptance

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/dns/v2/quotas"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestSchedulerStatsList(t *testing.T) {
	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	randomUUID := "513788b0-4f1b-4107-aee2-fbcdca9b9833"

	quotas, err := quotas.Get(client, randomUUID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotas)
}
