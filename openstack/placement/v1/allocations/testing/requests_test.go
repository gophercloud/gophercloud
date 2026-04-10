package testing

import (
	"context"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/placement/v1/allocations"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestGetSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetAllocationsSuccess(t, fakeServer)

	actual, err := allocations.Get(context.TODO(), client.ServiceClient(fakeServer), ConsumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocations, *actual)
}

func TestGetEmptySuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetEmptyAllocationsSuccess(t, fakeServer)

	actual, err := allocations.Get(context.TODO(), client.ServiceClient(fakeServer), EmptyConsumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedEmptyAllocations, *actual)
}

func TestUpdateSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdateAndGetAllocationsSuccess(t, fakeServer)

	err := allocations.Update(context.TODO(), client.ServiceClient(fakeServer), ConsumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			ProviderUUID1: {
				Resources: map[string]int{"VCPU": 2, "MEMORY_MB": 2048},
			},
		},
		ProjectID:          ProjectID,
		UserID:             UserID,
		ConsumerGeneration: nil,
	}).ExtractErr()
	th.AssertNoErr(t, err)

	actual, err := allocations.Get(context.TODO(), client.ServiceClient(fakeServer), ConsumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocationsAfterUpdate, *actual)
}

func TestUpdateNewConsumerSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdateAllocationsNewConsumerSuccess(t, fakeServer)

	// Act: Update with nil ConsumerGeneration; it must be serialized as JSON null, not omitted.
	err := allocations.Update(context.TODO(), client.ServiceClient(fakeServer), EmptyConsumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			ProviderUUID1: {
				Resources: map[string]int{"VCPU": 2, "MEMORY_MB": 2048},
			},
		},
		ProjectID:          ProjectID,
		UserID:             UserID,
		ConsumerGeneration: nil,
	}).ExtractErr()
	th.AssertNoErr(t, err)

	actual, err := allocations.Get(context.TODO(), client.ServiceClient(fakeServer), EmptyConsumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedAllocationsAfterUpdate, *actual)
}

func TestUpdateConflict(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleUpdateAllocationsConflict(t, fakeServer)

	staleGeneration := 0
	err := allocations.Update(context.TODO(), client.ServiceClient(fakeServer), ConflictConsumerUUID, allocations.UpdateOpts{
		Allocations: map[string]allocations.ProviderAllocationsOpts{
			ProviderUUID1: {
				Resources: map[string]int{"VCPU": 1},
			},
		},
		ProjectID:          ProjectID,
		UserID:             UserID,
		ConsumerGeneration: &staleGeneration,
	}).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}

func TestDeleteSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDeleteAndGetAllocationsSuccess(t, fakeServer)

	err := allocations.Delete(context.TODO(), client.ServiceClient(fakeServer), ConsumerUUID).ExtractErr()
	th.AssertNoErr(t, err)

	// Assert: Consumer now has no allocations.
	actual, err := allocations.Get(context.TODO(), client.ServiceClient(fakeServer), ConsumerUUID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, ExpectedEmptyAllocations, *actual)
}

func TestDeleteNotFound(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDeleteAllocationsNotFound(t, fakeServer)

	err := allocations.Delete(context.TODO(), client.ServiceClient(fakeServer), NotFoundConsumerUUID).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusNotFound))
}

func TestManageSuccess(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleManageAllocationsSuccess(t, fakeServer)

	err := allocations.Manage(context.TODO(), client.ServiceClient(fakeServer), allocations.ManageOpts{
		ManageConsumerUUID1: {
			Allocations: map[string]allocations.ProviderAllocationsOpts{
				ProviderUUID1: {
					Resources: map[string]int{"VCPU": 1},
				},
			},
			ProjectID:          ProjectID,
			UserID:             UserID,
			ConsumerGeneration: nil,
		},
		ManageConsumerUUID2: {
			Allocations: map[string]allocations.ProviderAllocationsOpts{
				ProviderUUID1: {
					Resources: map[string]int{"VCPU": 1},
				},
			},
			ProjectID:          ProjectID,
			UserID:             UserID,
			ConsumerGeneration: nil,
		},
	}).ExtractErr()
	th.AssertNoErr(t, err)

	for _, uuid := range []string{ManageConsumerUUID1, ManageConsumerUUID2} {
		actual, err := allocations.Get(context.TODO(), client.ServiceClient(fakeServer), uuid).Extract()
		th.AssertNoErr(t, err)
		th.AssertDeepEquals(t, ExpectedAllocationsAfterManage, *actual)
	}
}

func TestManageConflict(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleManageAllocationsConflict(t, fakeServer)

	staleGeneration := 0
	err := allocations.Manage(context.TODO(), client.ServiceClient(fakeServer), allocations.ManageOpts{
		ConflictConsumerUUID: {
			Allocations: map[string]allocations.ProviderAllocationsOpts{
				ProviderUUID1: {
					Resources: map[string]int{"VCPU": 1},
				},
			},
			ProjectID:          ProjectID,
			UserID:             UserID,
			ConsumerGeneration: &staleGeneration,
		},
	}).ExtractErr()
	th.AssertEquals(t, true, gophercloud.ResponseCodeIs(err, http.StatusConflict))
}
