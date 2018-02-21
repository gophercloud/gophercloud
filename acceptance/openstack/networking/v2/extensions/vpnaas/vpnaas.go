package vpnaas

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/ipsecpolicies"
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

// CreateIPSecPolicy will create an IPSec Policy with a random name and given
// rule. An error will be returned if the rule could not be created.
func CreateIPSecPolicy(t *testing.T, client *gophercloud.ServiceClient) (*ipsecpolicies.Policy, error) {
	policyName := tools.RandomString("TESTACC-", 8)

	t.Logf("Attempting to create policy %s", policyName)

	createOpts := ipsecpolicies.CreateOpts{
		Name: policyName,
	}

	policy, err := ipsecpolicies.Create(client, createOpts).Extract()
	if err != nil {
		return policy, err
	}

	t.Logf("Successfully created IPSec policy %s", policyName)

	return policy, nil
}
