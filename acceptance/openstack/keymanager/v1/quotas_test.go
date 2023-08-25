//go:build acceptance || keymanager || orders
// +build acceptance keymanager orders

package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/quotas"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestQuotasCRUD(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewKeyManagerV1Client()
	th.AssertNoErr(t, err)

	// Get a project
	identityClient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)
	enabled := true
	projectPages, err := projects.List(identityClient, projects.ListOpts{Enabled: &enabled}).AllPages()
	th.AssertNoErr(t, err)

	projectList, err := projects.ExtractProjects(projectPages)
	th.AssertNoErr(t, err)

	project := projectList[0]

	// Set a quota
	err = quotas.Update(client, project.ID, quotas.UpdateOpts{
		Secrets:    gophercloud.IntToPointer(1),
		Orders:     gophercloud.IntToPointer(2),
		Containers: gophercloud.IntToPointer(3),
		Consumers:  gophercloud.IntToPointer(4),
		CAS:        gophercloud.IntToPointer(5),
	}).Err
	th.AssertNoErr(t, err)

	// Get quota
	quota, err := quotas.GetProjectQuota(client, project.ID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &quotas.Quota{
		Secrets:    gophercloud.IntToPointer(1),
		Orders:     gophercloud.IntToPointer(2),
		Containers: gophercloud.IntToPointer(3),
		Consumers:  gophercloud.IntToPointer(4),
		CAS:        gophercloud.IntToPointer(5),
	}, quota)

	// Delete quota
	err = quotas.Delete(client, project.ID).Err
	th.AssertNoErr(t, err)
}
