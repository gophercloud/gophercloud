//go:build acceptance
// +build acceptance

package v3

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/identity/v3/extensions/federation"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestListMappings(t *testing.T) {
	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	allPages, err := federation.ListMappings(client).AllPages()
	th.AssertNoErr(t, err)

	mappings, err := federation.ExtractMappings(allPages)
	th.AssertNoErr(t, err)

	tools.PrintResource(t, mappings)
}
