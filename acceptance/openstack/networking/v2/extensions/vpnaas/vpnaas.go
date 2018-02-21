package vpnaas

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"
)

// CreateService will create a Service with a random name and a specified router ID
// An error will be returned if the service could not be created.
func CreateService(t *testing.T, client *gophercloud.ServiceClient, routerID string) (*services.Service, error) {
	serviceName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create service %s", serviceName)

	iTrue := true
	createOpts := services.CreateOpts{
		Name:         serviceName,
		AdminStateUp: &iTrue,
		RouterID:     routerID,
	}
	service, err := services.Create(client, createOpts).Extract()
	if err != nil {
		return service, err
	}

	t.Logf("Successfully created service %s", serviceName)

	return service, nil
}

// DeleteService will delete a service with a specified ID. A fatal error
// will occur if the delete was not successful. This works best when used as
// a deferred function.
func DeleteService(t *testing.T, client *gophercloud.ServiceClient, serviceID string) {
	t.Logf("Attempting to delete service: %s", serviceID)

	err := services.Delete(client, serviceID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete service %s: %v", serviceID, err)
	}

	t.Logf("Service deleted: %s", serviceID)
}
