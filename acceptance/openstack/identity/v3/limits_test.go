//go:build acceptance
// +build acceptance

package v3

import (
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/limits"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestGetEnforcementModel(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	model, err := limits.GetEnforcementModel(client).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, model)
}

func TestLimitsList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	listOpts := limits.ListOpts{}

	allPages, err := limits.List(client, listOpts).AllPages()
	th.AssertNoErr(t, err)

	_, err = limits.ExtractLimits(allPages)
	th.AssertNoErr(t, err)
}

func TestLimitsCRUD(t *testing.T) {
	// TODO: After https://github.com/gophercloud/gophercloud/issues/2290 is implemented
	// use registered limits API to create global registered limit and then overwrite it with limit.
	// Current solution (using glance limit) only works with Openstack Xena and above.
	clients.SkipReleasesBelow(t, "stable/xena")

	err := os.Setenv("OS_SYSTEM_SCOPE", "all")
	th.AssertNoErr(t, err)
	defer os.Unsetenv("OS_SYSTEM_SCOPE")

	limitDescription := tools.RandomString("TESTLIMITS-DESC-", 8)
	resourceLimit := tools.RandomInt(1, 100)
	resourceName := "image_size_total"

	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)

	// Find image service (glance on Devstack) which has precreated registered limits.
	allPages, err := services.List(client, nil).AllPages()
	th.AssertNoErr(t, err)

	svList, err := services.ExtractServices(allPages)
	serviceID := ""
	for _, service := range svList {
		if service.Type == "image" {
			serviceID = service.ID
			break
		}
	}
	th.AssertIntGreaterOrEqual(t, len(serviceID), 1)

	createOpts := limits.BatchCreateOpts{
		limits.CreateOpts{
			ServiceID:     serviceID,
			ProjectID:     project.ID,
			ResourceName:  resourceName,
			ResourceLimit: resourceLimit,
			Description:   limitDescription,
			RegionID:      "RegionOne",
		},
	}

	createdLimits, err := limits.BatchCreate(client, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertIntGreaterOrEqual(t, 1, len(createdLimits))
	th.AssertEquals(t, limitDescription, createdLimits[0].Description)
	th.AssertEquals(t, resourceLimit, createdLimits[0].ResourceLimit)
	th.AssertEquals(t, resourceName, createdLimits[0].ResourceName)
	th.AssertEquals(t, serviceID, createdLimits[0].ServiceID)
	th.AssertEquals(t, project.ID, createdLimits[0].ProjectID)

	limitID := createdLimits[0].ID

	limit, err := limits.Get(client, limitID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, createdLimits[0], *limit)

	newLimitDescription := tools.RandomString("TESTLIMITS-DESC-CHNGD-", 8)
	newResourceLimit := tools.RandomInt(1, 100)
	updateOpts := limits.UpdateOpts{
		Description:   &newLimitDescription,
		ResourceLimit: &newResourceLimit,
	}

	updatedLimit, err := limits.Update(client, limitID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newLimitDescription, updatedLimit.Description)
	th.AssertEquals(t, newResourceLimit, updatedLimit.ResourceLimit)
}
