package v1

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/placement/v1/resourceproviders"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func CreateResourceProvider(t *testing.T, client *gophercloud.ServiceClient) (*resourceproviders.ResourceProvider, error) {
	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create resource provider: %s", name)

	createOpts := resourceproviders.CreateOpts{
		Name: name,
	}

	client.Microversion = "1.20"
	resourceProvider, err := resourceproviders.Create(client, createOpts).Extract()
	if err != nil {
		return resourceProvider, err
	}

	t.Logf("Successfully created resourceProvider: %s.", resourceProvider.Name)
	tools.PrintResource(t, resourceProvider)

	th.AssertEquals(t, resourceProvider.Name, name)

	return resourceProvider, nil
}

func CreateResourceProviderWithParent(t *testing.T, client *gophercloud.ServiceClient, parentUUID string) (*resourceproviders.ResourceProvider, error) {
	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create resource provider: %s", name)

	createOpts := resourceproviders.CreateOpts{
		Name:               name,
		ParentProviderUUID: parentUUID,
	}

	client.Microversion = "1.20"
	resourceProvider, err := resourceproviders.Create(client, createOpts).Extract()
	if err != nil {
		return resourceProvider, err
	}

	t.Logf("Successfully created resourceProvider: %s.", resourceProvider.Name)
	tools.PrintResource(t, resourceProvider)

	th.AssertEquals(t, resourceProvider.Name, name)
	th.AssertEquals(t, resourceProvider.ParentProviderUUID, parentUUID)

	return resourceProvider, nil
}
