//go:build acceptance || identity || endpointgroups
// +build acceptance identity endpointgroups

package v3

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/endpointgroups"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestEndpointGroupsList(t *testing.T) {
	clients.SkipRelease(t, "stable/mitaka")
	clients.SkipRelease(t, "stable/newton")
	clients.SkipRelease(t, "stable/ocata")
	clients.SkipRelease(t, "stable/pike")
	clients.SkipRelease(t, "stable/queens")
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := endpointgroups.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allEndpointGroups, err := endpointgroups.ExtractEndpointGroups(allPages)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, len(allEndpointGroups), 0)
}
