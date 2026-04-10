//go:build acceptance || placement || allocations

package v1

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/allocations"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestGetAllocationsSuccess(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	consumerUUID := fmt.Sprintf("%08x-0000-0000-0000-000000000000", rand.Int31())

	// Assert: We don't have any allocations for this random UUID.
	// We get an empty allocations map, not 404.
	allocs, err := allocations.Get(context.TODO(), client, consumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(allocs.Allocations))
}

func TestUpdateAllocationsNewConsumerSuccess(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/rocky")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, _, err := CreateResourceProviderWithVCPUInventory(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	consumerUUID := fmt.Sprintf("%08x-0000-0000-0000-000000000000", rand.Int31())
	defer allocations.Delete(context.TODO(), client, consumerUUID)

	client.Microversion = "1.28"

	// Act: Update with nil ConsumerGeneration to signal a new consumer (serialized as JSON null, not omitted).
	err = allocations.Update(context.TODO(), client, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			resourceProvider.UUID: {
				Resources: map[string]int{"VCPU": 2, "MEMORY_MB": 1024},
			},
		},
		ProjectID:          "test-project",
		UserID:             "test-user",
		ConsumerGeneration: nil,
	}).ExtractErr()
	th.AssertNoErr(t, err)

	allocs, err := allocations.Get(context.TODO(), client, consumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, len(allocs.Allocations))
	th.AssertEquals(t, 2, allocs.Allocations[resourceProvider.UUID].Resources["VCPU"])
	th.AssertEquals(t, 1024, allocs.Allocations[resourceProvider.UUID].Resources["MEMORY_MB"])
}

func TestUpdateAllocationsSuccess(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/rocky")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, _, err := CreateResourceProviderWithVCPUInventory(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	consumerUUID := fmt.Sprintf("%08x-0000-0000-0000-000000000001", rand.Int31())
	defer allocations.Delete(context.TODO(), client, consumerUUID)

	client.Microversion = "1.28"

	// Arrange: Create the consumer with nil ConsumerGeneration.
	err = allocations.Update(context.TODO(), client, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			resourceProvider.UUID: {
				Resources: map[string]int{"VCPU": 1},
			},
		},
		ProjectID:          "test-project",
		UserID:             "test-user",
		ConsumerGeneration: nil,
	}).ExtractErr()
	th.AssertNoErr(t, err)

	existing, err := allocations.Get(context.TODO(), client, consumerUUID).Extract()
	th.AssertNoErr(t, err)

	// Act: Update allocations using the consumer's current generation.
	err = allocations.Update(context.TODO(), client, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			resourceProvider.UUID: {
				Resources: map[string]int{"VCPU": 2, "MEMORY_MB": 1024},
			},
		},
		ProjectID:          *existing.ProjectID,
		UserID:             *existing.UserID,
		ConsumerGeneration: existing.ConsumerGeneration,
	}).ExtractErr()
	th.AssertNoErr(t, err)

	updated, err := allocations.Get(context.TODO(), client, consumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, updated.Allocations[resourceProvider.UUID].Resources["VCPU"])
	th.AssertEquals(t, 1024, updated.Allocations[resourceProvider.UUID].Resources["MEMORY_MB"])

	tools.PrintResource(t, updated)
}

func TestUpdateAllocationsConflict(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/rocky")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, _, err := CreateResourceProviderWithVCPUInventory(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	consumerUUID := fmt.Sprintf("%08x-0000-0000-0000-000000000002", rand.Int31())
	defer allocations.Delete(context.TODO(), client, consumerUUID)

	client.Microversion = "1.28"

	// Arrange: Create the consumer to establish a valid generation.
	err = allocations.Update(context.TODO(), client, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			resourceProvider.UUID: {
				Resources: map[string]int{"VCPU": 1},
			},
		},
		ProjectID:          "test-project",
		UserID:             "test-user",
		ConsumerGeneration: nil,
	}).ExtractErr()
	th.AssertNoErr(t, err)

	// Act: Update with a stale generation to trigger a 409 conflict.
	staleGeneration := -1
	err = allocations.Update(context.TODO(), client, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			resourceProvider.UUID: {
				Resources: map[string]int{"VCPU": 2},
			},
		},
		ProjectID:          "test-project",
		UserID:             "test-user",
		ConsumerGeneration: &staleGeneration,
	}).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

func TestDeleteAllocationsSuccess(t *testing.T) {
	clients.RequireAdmin(t)
	clients.SkipReleasesBelow(t, "stable/rocky")

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, _, err := CreateResourceProviderWithVCPUInventory(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	consumerUUID := fmt.Sprintf("%08x-0000-0000-0000-000000000003", rand.Int31())

	client.Microversion = "1.28"

	// Arrange: Create allocations for the consumer.
	err = allocations.Update(context.TODO(), client, consumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			resourceProvider.UUID: {
				Resources: map[string]int{"VCPU": 1},
			},
		},
		ProjectID:          "test-project",
		UserID:             "test-user",
		ConsumerGeneration: nil,
	}).ExtractErr()
	th.AssertNoErr(t, err)

	// Act: Delete all allocations for the consumer.
	err = allocations.Delete(context.TODO(), client, consumerUUID).ExtractErr()
	th.AssertNoErr(t, err)

	// Assert: Consumer now returns an empty allocations map.
	allocs, err := allocations.Get(context.TODO(), client, consumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 0, len(allocs.Allocations))
}

func TestDeleteAllocationsNotFound(t *testing.T) {
	clients.RequireAdmin(t)

	client, err := clients.NewPlacementV1Client()
	th.AssertNoErr(t, err)

	resourceProvider, err := CreateResourceProvider(t, client)
	th.AssertNoErr(t, err)
	defer DeleteResourceProvider(t, client, resourceProvider.UUID)

	// Assert: An RP that was never a consumer returns 404 on DELETE.
	err = allocations.Delete(context.TODO(), client, resourceProvider.UUID).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}
