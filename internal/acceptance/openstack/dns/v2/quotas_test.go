//go:build acceptance || dns || quotas

package v2

import (
	"context"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/quotas"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestQuotaGetUpdate(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewDNSV2Client()
	th.AssertNoErr(t, err)

	identityClient, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	projectID, err := getProjectID(t, identityClient)
	th.AssertNoErr(t, err)

	// test Get Quota
	quota, err := quotas.Get(context.TODO(), client, projectID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quota)

	// test Update Quota
	zones := 9
	updateOpts := quotas.UpdateOpts{
		Zones: &zones,
	}
	res, err := quotas.Update(context.TODO(), client, projectID, updateOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, res)

	quota.Zones = zones
	th.AssertDeepEquals(t, *quota, *res)
}

func getProjectID(t *testing.T, client *gophercloud.ServiceClient) (string, error) {
	allPages, err := projects.ListAvailable(client).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allProjects, err := projects.ExtractProjects(allPages)
	th.AssertNoErr(t, err)

	for _, project := range allProjects {
		return project.ID, nil
	}

	return "", fmt.Errorf("Unable to get project ID")
}
