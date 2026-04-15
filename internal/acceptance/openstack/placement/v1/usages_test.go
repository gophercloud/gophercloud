//go:build acceptance || placement || usages

package v1

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/allocations"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/usages"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGetUsagesSuccess(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/rocky")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, _, err := CreateResourceProviderWithVCPUInventory(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	consumerUUID := fmt.Sprintf("%08x-0000-0000-0000-000000000010", rand.Int31())
	defer allocations.Delete(context.TODO(), client, consumerUUID)

	projectID := "test-project"
	userID := "test-user"

	client.Microversion = "1.38"

	// Arrange: Create allocations for a consumer under the test project.
	err = allocations.Update(context.TODO(), client, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			resourceProvider.UUID: {
				Resources: map[string]int{"VCPU": 2, "MEMORY_MB": 1024},
			},
		},
		ProjectID:          projectID,
		UserID:             userID,
		ConsumerGeneration: nil,
		ConsumerType:       "INSTANCE",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	// Act: Retrieve total usages for the project (1.38+, grouped by consumer type).
	totalUsages, err := usages.Get(context.TODO(), client, usages.GetOpts{
		ProjectID: projectID,
	}).Extract()
	th.AssertNoErr(t, err)

	// Assert: Exactly one consumer type "INSTANCE" with our exact resource usage.
	th.AssertEquals(t, 1, len(totalUsages.Usages))
	ctUsage := totalUsages.Usages["INSTANCE"]
	th.AssertEquals(t, 2, ctUsage["VCPU"])
	th.AssertEquals(t, 1024, ctUsage["MEMORY_MB"])
	th.AssertEquals(t, 1, ctUsage["consumer_count"])

	tools.PrintResource(t, totalUsages)
}

func TestGetUsagesWithUserSuccess(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/rocky")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, _, err := CreateResourceProviderWithVCPUInventory(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	consumerUUID := fmt.Sprintf("%08x-0000-0000-0000-000000000011", rand.Int31())
	defer allocations.Delete(context.TODO(), client, consumerUUID)

	projectID := "test-project"
	userID := "test-user"

	client.Microversion = "1.38"

	// Arrange: Create allocations.
	err = allocations.Update(context.TODO(), client, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			resourceProvider.UUID: {
				Resources: map[string]int{"VCPU": 1},
			},
		},
		ProjectID:          projectID,
		UserID:             userID,
		ConsumerGeneration: nil,
		ConsumerType:       "INSTANCE",
	}).ExtractErr()
	th.AssertNoErr(t, err)

	// Act: Retrieve usages filtered by project and user (1.38+).
	totalUsages, err := usages.Get(context.TODO(), client, usages.GetOpts{
		ProjectID: projectID,
		UserID:    userID,
	}).Extract()
	th.AssertNoErr(t, err)

	// Assert: Exactly one consumer type "INSTANCE" with our exact VCPU usage.
	th.AssertEquals(t, 1, len(totalUsages.Usages))
	ctUsage := totalUsages.Usages["INSTANCE"]
	th.AssertEquals(t, 1, ctUsage["VCPU"])
	th.AssertEquals(t, 1, ctUsage["consumer_count"])

	tools.PrintResource(t, totalUsages)
}

func TestGetUsagesEmptySuccess(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/rocky")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.38"

	// Act: Query usages for a project with no allocations.
	totalUsages, err := usages.Get(context.TODO(), client, usages.GetOpts{
		ProjectID: "nonexistent-project-with-no-allocations",
	}).Extract()
	th.AssertNoErr(t, err)

	// Assert: Empty usages map.
	th.AssertEquals(t, 0, len(totalUsages.Usages))
}

func TestGetUsagesPre138Success(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/rocky")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, _, err := CreateResourceProviderWithVCPUInventory(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	consumerUUID := fmt.Sprintf("%08x-0000-0000-0000-000000000012", rand.Int31())
	defer allocations.Delete(context.TODO(), client, consumerUUID)

	projectID := "test-project"
	userID := "test-user"

	client.Microversion = "1.28"

	// Arrange: Create allocations for a consumer.
	err = allocations.Update(context.TODO(), client, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			resourceProvider.UUID: {
				Resources: map[string]int{"VCPU": 2, "MEMORY_MB": 1024},
			},
		},
		ProjectID:          projectID,
		UserID:             userID,
		ConsumerGeneration: nil,
	}).ExtractErr()
	th.AssertNoErr(t, err)

	// Act: Retrieve total usages for the project (pre-1.38, flat map).
	totalUsages, err := usages.Get(context.TODO(), client, usages.GetOpts{
		ProjectID: projectID,
	}).ExtractPre138()
	th.AssertNoErr(t, err)

	// Assert: Usages reflect the allocations we just created.
	th.AssertEquals(t, 2, totalUsages.Usages["VCPU"])
	th.AssertEquals(t, 1024, totalUsages.Usages["MEMORY_MB"])

	tools.PrintResource(t, totalUsages)
}

func TestGetUsagesPre138EmptySuccess(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/rocky")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	client.Microversion = "1.9"

	// Act: Query usages for a project with no allocations (pre-1.38).
	totalUsages, err := usages.Get(context.TODO(), client, usages.GetOpts{
		ProjectID: "nonexistent-project-with-no-allocations",
	}).ExtractPre138()
	th.AssertNoErr(t, err)

	// Assert: Empty usages map.
	th.AssertEquals(t, 0, len(totalUsages.Usages))
}
