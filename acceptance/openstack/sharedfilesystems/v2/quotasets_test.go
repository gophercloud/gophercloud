//go:build acceptance
// +build acceptance

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestGet(t *testing.T) {
	client, err := clients.NewSharedFilesystemV2Client()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	quotaset, err := quotasets.Get(client, client.TenantID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestUpdate(t *testing.T) {
	client, err := clients.NewSharedFilesystemV2Client()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	quotaset, err := quotasets.Get(client, client.TenantID).Extract()
	th.AssertNoErr(t, err)

	// Update the quotaset
	updateOpts := quotasets.UpdateOpts{
		Gigabytes: gophercloud.Int(100),
	}
	quotaset, err = quotasets.Update(client, client.TenantID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestGetByShareType(t *testing.T) {
	client, err := clients.NewSharedFilesystemV2Client()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	quotaset, err := quotasets.GetByShareType(client, client.TenantID, "default").Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestUpdateByShareType(t *testing.T) {
	client, err := clients.NewSharedFilesystemV2Client()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	quotaset, err := quotasets.GetByShareType(client, client.TenantID, "default").Extract()
	th.AssertNoErr(t, err)

	// Update the quotaset
	updateOpts := quotasets.UpdateOpts{
		Gigabytes: gophercloud.Int(100),
	}
	quotaset, err = quotasets.UpdateByShareType(client, client.TenantID, "default", updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestGetByUser(t *testing.T) {
	client, err := clients.NewSharedFilesystemV2Client()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	quotaset, err := quotasets.GetByUser(client, client.TenantID, "admin").Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}

func TestUpdateByUser(t *testing.T) {
	client, err := clients.NewSharedFilesystemV2Client()
	th.AssertNoErr(t, err)

	// Get the quotaset for the current tenant
	quotaset, err := quotasets.GetByUser(client, client.TenantID, "admin").Extract()
	th.AssertNoErr(t, err)

	// Update the quotaset
	updateOpts := quotasets.UpdateOpts{
		Gigabytes: gophercloud.Int(100),
	}
	quotaset, err = quotasets.UpdateByUser(client, client.TenantID, "admin", updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaset)
}
