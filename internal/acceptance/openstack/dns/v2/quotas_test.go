//go:build acceptance || dns || quotas

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	identity "github.com/gophercloud/gophercloud/v2/internal/acceptance/openstack/identity/v3"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/quotas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestQuotaGetUpdate(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	identityClient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := identity.CreateProject(t, identityClient, nil)
	th.AssertNoErr(t, err)
	defer identity.DeleteProject(t, identityClient, project.ID)

	// use DNS specific header to set the project ID
	client.MoreHeaders = map[string]string{
		"X-Auth-Sudo-Tenant-ID": project.ID,
	}

	// test Get Quota
	quota, err := quotas.Get(context.TODO(), client, project.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quota)

	// test Update Quota
	zones := 9
	updateOpts := quotas.UpdateOpts{
		Zones: &zones,
	}
	res, err := quotas.Update(context.TODO(), client, project.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, res)

	quota.Zones = zones
	th.AssertDeepEquals(t, *quota, *res)
}
