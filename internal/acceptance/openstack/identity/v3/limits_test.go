//go:build acceptance || identity || limits

package v3

import (
	"context"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/limits"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/registeredlimits"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/services"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGetEnforcementModel(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	model, err := limits.GetEnforcementModel(context.TODO(), client).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, model)
}

func TestLimitsList(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	listOpts := limits.ListOpts{}

	allPages, err := limits.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	_, err = limits.ExtractLimits(allPages)
	th.AssertNoErr(t, err)
}

func TestLimitsCRUD(t *testing.T) {
	err := os.Setenv("OS_SYSTEM_SCOPE", "all")
	th.AssertNoErr(t, err)
	defer os.Unsetenv("OS_SYSTEM_SCOPE")

	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)

	// Get the service to register the limit against.
	allPages, err := services.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	svList, err := services.ExtractServices(allPages)
	th.AssertNoErr(t, err)
	serviceID := ""
	for _, service := range svList {
		serviceID = service.ID
		break
	}
	th.AssertIntGreaterOrEqual(t, len(serviceID), 1)

	// Create global registered limit
	description := tools.RandomString("GLOBALLIMIT-DESC-", 8)
	defaultLimit := tools.RandomInt(1, 100)
	globalResourceName := tools.RandomString("GLOBALLIMIT-", 8)

	createRegisteredLimitsOpts := registeredlimits.BatchCreateOpts{
		registeredlimits.CreateOpts{
			ServiceID:    serviceID,
			ResourceName: globalResourceName,
			DefaultLimit: defaultLimit,
			Description:  description,
			RegionID:     "RegionOne",
		},
	}

	createdRegisteredLimits, err := registeredlimits.BatchCreate(context.TODO(), client, createRegisteredLimitsOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, createdRegisteredLimits[0])
	th.AssertIntGreaterOrEqual(t, 1, len(createdRegisteredLimits))

	// Override global limit in specific project
	limitDescription := tools.RandomString("TESTLIMITS-DESC-", 8)
	resourceLimit := tools.RandomInt(1, 1000)

	createOpts := limits.BatchCreateOpts{
		limits.CreateOpts{
			ServiceID:     serviceID,
			ProjectID:     project.ID,
			ResourceName:  globalResourceName,
			ResourceLimit: resourceLimit,
			Description:   limitDescription,
			RegionID:      "RegionOne",
		},
	}

	createdLimits, err := limits.BatchCreate(context.TODO(), client, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertIntGreaterOrEqual(t, 1, len(createdLimits))
	th.AssertEquals(t, limitDescription, createdLimits[0].Description)
	th.AssertEquals(t, resourceLimit, createdLimits[0].ResourceLimit)
	th.AssertEquals(t, globalResourceName, createdLimits[0].ResourceName)
	th.AssertEquals(t, serviceID, createdLimits[0].ServiceID)
	th.AssertEquals(t, project.ID, createdLimits[0].ProjectID)

	limitID := createdLimits[0].ID

	limit, err := limits.Get(context.TODO(), client, limitID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, createdLimits[0], *limit)

	newLimitDescription := tools.RandomString("TESTLIMITS-DESC-CHNGD-", 8)
	newResourceLimit := tools.RandomInt(1, 100)
	updateOpts := limits.UpdateOpts{
		Description:   &newLimitDescription,
		ResourceLimit: &newResourceLimit,
	}

	updatedLimit, err := limits.Update(context.TODO(), client, limitID, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, newLimitDescription, updatedLimit.Description)
	th.AssertEquals(t, newResourceLimit, updatedLimit.ResourceLimit)

	// Verify Deleting registered limit fails as it has project specific limit associated with it
	del_err := registeredlimits.Delete(context.TODO(), client, createdRegisteredLimits[0].ID).ExtractErr()
	th.AssertErr(t, del_err)

	// Delete project specific limit
	err = limits.Delete(context.TODO(), client, limitID).ExtractErr()
	th.AssertNoErr(t, err)

	_, err = limits.Get(context.TODO(), client, limitID).Extract()
	th.AssertErr(t, err)

	// Delete registered limit
	err = registeredlimits.Delete(context.TODO(), client, createdRegisteredLimits[0].ID).ExtractErr()
	th.AssertNoErr(t, err)
}
