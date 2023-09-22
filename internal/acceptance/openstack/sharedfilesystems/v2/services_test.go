//go:build acceptance
// +build acceptance

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/services"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestServicesList(t *testing.T) {
	client, err := clients.NewSharedFileSystemV2Client()
	th.AssertNoErr(t, err)

	client.Microversion = "2.7"
	allPages, err := services.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allServices, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)

	th.AssertIntGreaterOrEqual(t, len(allServices), 1)

	for _, s := range allServices {
		tools.PrintResource(t, &s)
		th.AssertEquals(t, s.Status, "enabled")
	}
}
