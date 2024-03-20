//go:build acceptance || identity || tenants

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/tenants"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestTenantsList(t *testing.T) {
	clients.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2Client()
	th.AssertNoErr(t, err)

	allPages, err := tenants.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allTenants, err := tenants.ExtractTenants(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, tenant := range allTenants {
		tools.PrintResource(t, tenant)

		if tenant.Name == "admin" {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestTenantsCRUD(t *testing.T) {
	clients.RequireIdentityV2(t)
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV2AdminClient()
	th.AssertNoErr(t, err)

	tenant, err := CreateTenant(t, client, nil)
	th.AssertNoErr(t, err)
	defer DeleteTenant(t, client, tenant.ID)

	tenant, err = tenants.Get(context.TODO(), client, tenant.ID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, tenant)

	description := ""
	updateOpts := tenants.UpdateOpts{
		Description: &description,
	}

	newTenant, err := tenants.Update(context.TODO(), client, tenant.ID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, newTenant)

	th.AssertEquals(t, newTenant.Description, description)
}
