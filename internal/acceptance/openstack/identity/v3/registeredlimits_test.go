//go:build acceptance || identity || registeredlimits

package v3

import (
	"context"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/registeredlimits"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/services"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestRegisteredLimitsCRUD(t *testing.T) {
	err := os.Setenv("OS_SYSTEM_SCOPE", "all")
	th.AssertNoErr(t, err)
	defer os.Unsetenv("OS_SYSTEM_SCOPE")

	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	// Get glance service to register the limit
	allServicePages, err := services.List(client, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	svList, err := services.ExtractServices(allServicePages)
	th.AssertNoErr(t, err)
	serviceID := ""
	for _, service := range svList {
		serviceID = service.ID
		break
	}
	th.AssertIntGreaterOrEqual(t, len(serviceID), 1)

	// Create RegisteredLimit
	limitDescription := tools.RandomString("TESTLIMITS-DESC-", 8)
	defaultLimit := tools.RandomInt(1, 100)
	resourceName := tools.RandomString("LIMIT-NAME-", 8)

	createOpts := registeredlimits.BatchCreateOpts{
		registeredlimits.CreateOpts{
			ServiceID:    serviceID,
			ResourceName: resourceName,
			DefaultLimit: defaultLimit,
			Description:  limitDescription,
			RegionID:     "RegionOne",
		},
	}

	createdRegisteredLimits, err := registeredlimits.BatchCreate(context.TODO(), client, createOpts).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, createdRegisteredLimits[0])
	th.AssertIntGreaterOrEqual(t, 1, len(createdRegisteredLimits))
	th.AssertEquals(t, limitDescription, createdRegisteredLimits[0].Description)
	th.AssertEquals(t, defaultLimit, createdRegisteredLimits[0].DefaultLimit)
	th.AssertEquals(t, resourceName, createdRegisteredLimits[0].ResourceName)
	th.AssertEquals(t, serviceID, createdRegisteredLimits[0].ServiceID)
	th.AssertEquals(t, "RegionOne", createdRegisteredLimits[0].RegionID)

	// List the registered limits
	listOpts := registeredlimits.ListOpts{}
	allPages, err := registeredlimits.List(client, listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	_, err = registeredlimits.ExtractRegisteredLimits(allPages)
	th.AssertNoErr(t, err)

	// Get RegisteredLimit by ID
	registered_limit, err := registeredlimits.Get(context.TODO(), client, createdRegisteredLimits[0].ID).Extract()
	th.AssertNoErr(t, err)
	tools.PrintResource(t, registered_limit)

	// Update the existing registered_limit
	updatedDescription := "Test description for registered limit"
	updatedDefaultLimit := 1000
	updatedResourceName := tools.RandomString("LIMIT-NAME-", 8)
	updatedOpts := registeredlimits.UpdateOpts{
		Description:  &updatedDescription,
		DefaultLimit: &updatedDefaultLimit,
		ServiceID:    serviceID,
		ResourceName: updatedResourceName,
	}

	updated_registered_limit, err := registeredlimits.Update(context.TODO(), client, createdRegisteredLimits[0].ID, updatedOpts).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, updated_registered_limit)
	th.AssertEquals(t, updated_registered_limit.Description, updatedDescription)
	th.AssertEquals(t, updated_registered_limit.DefaultLimit, updatedDefaultLimit)
	th.AssertEquals(t, updated_registered_limit.ResourceName, updatedResourceName)

	// Delete the registered limit
	del_err := registeredlimits.Delete(context.TODO(), client, createdRegisteredLimits[0].ID).ExtractErr()
	th.AssertNoErr(t, del_err)

	_, err = registeredlimits.Get(context.TODO(), client, createdRegisteredLimits[0].ID).Extract()
	th.AssertErr(t, err)
}
