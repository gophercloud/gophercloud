package v1

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/resourceproviders"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func CreateResourceProvider(t *testing.T, client *gophercloud.ServiceClient) (*resourceproviders.ResourceProvider, error) {
	name := tools.RandomString("TESTACC-", 8)
	t.Logf("Attempting to create resource provider: %s", name)

	createOpts := resourceproviders.CreateOpts{
		Name: name,
	}

	client.Microversion = "1.20"
	resourceProvider, err := resourceproviders.Create(context.TODO(), client, createOpts).Extract()
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
	resourceProvider, err := resourceproviders.Create(context.TODO(), client, createOpts).Extract()
	if err != nil {
		return resourceProvider, err
	}

	t.Logf("Successfully created resourceProvider: %s.", resourceProvider.Name)
	tools.PrintResource(t, resourceProvider)

	th.AssertEquals(t, resourceProvider.Name, name)
	th.AssertEquals(t, resourceProvider.ParentProviderUUID, parentUUID)

	return resourceProvider, nil
}

// DeleteResourceProvider will delete a resource provider with a specified ID.
// A fatal error will occur if the delete was not successful. This works best when
// used as a deferred function.
func DeleteResourceProvider(t *testing.T, client *gophercloud.ServiceClient, resourceProviderID string) {
	t.Logf("Attempting to delete resourceProvider: %s", resourceProviderID)

	err := resourceproviders.Delete(context.TODO(), client, resourceProviderID).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete resourceProvider %s: %v", resourceProviderID, err)
	}

	t.Logf("Deleted resourceProvider: %s.", resourceProviderID)
}
