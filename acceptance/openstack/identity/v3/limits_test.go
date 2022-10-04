//go:build acceptance
// +build acceptance

package v3

import (
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

func TestCreateLimits(t *testing.T) {
	limitDescription := tools.RandomString("TESTLIMITS-DESC-", 8)
	resourceLimit := tools.RandomInt(1, 100)
	resourceName := "volume"

	clients.RequireAdmin(t)

	client, err := clients.NewIdentityV3Client()
	th.AssertNoErr(t, err)

	project, err := CreateProject(t, client, nil)
	th.AssertNoErr(t, err)

	createServiceOpts := &services.CreateOpts{
		Type:  resourceName,
		Extra: map[string]interface{}{},
	}

	service, err := CreateService(t, client, createServiceOpts)
	th.AssertNoErr(t, err)

	createOpts := limits.BatchCreateOpts{
		limits.CreateOpts{
			ServiceID:     service.ID,
			ProjectID:     project.ID,
			ResourceName:  resourceName,
			ResourceLimit: resourceLimit,
			Description:   limitDescription,
		},
	}

	createdLimits, err := limits.BatchCreate(client, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertIntGreaterOrEqual(t, 1, len(createdLimits))
	th.AssertEquals(t, limitDescription, createdLimits[0].Description)
	th.AssertEquals(t, resourceLimit, createdLimits[0].ResourceLimit)
	th.AssertEquals(t, resourceName, createdLimits[0].ResourceName)
	th.AssertEquals(t, service.ID, createdLimits[0].ServiceID)
	th.AssertEquals(t, project.ID, createdLimits[0].ProjectID)
}
