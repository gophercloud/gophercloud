package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/replicas"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func getClient(microVersion string) *gophercloud.ServiceClient {
	c := client.ServiceClient()
	c.Type = "sharev2"
	c.Microversion = microVersion
	return c
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := &replicas.CreateOpts{
		ShareID:          "65a34695-f9e5-4eea-b48d-a0b261d82943",
		AvailabilityZone: "zone-1",
	}
	actual, err := replicas.Create(context.TODO(), getClient("2.11"), options).Extract()

	expected := &replicas.Replica{
		ID:               "3b9c33e8-b136-45c6-84a6-019c8db1d550",
		ShareID:          "65a34695-f9e5-4eea-b48d-a0b261d82943",
		AvailabilityZone: "zone-1",
		Status:           "creating",
		ShareNetworkID:   "ca0163c8-3941-4420-8b01-41517e19e366",
		CreatedAt:        time.Date(2023, time.May, 26, 12, 32, 56, 391337000, time.UTC), //"2023-05-26T12:32:56.391337",
	}

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	result := replicas.Delete(context.TODO(), getClient("2.11"), replicaID)
	th.AssertNoErr(t, result.Err)
}

func TestForceDeleteSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockForceDeleteResponse(t)

	err := replicas.ForceDelete(context.TODO(), getClient("2.11"), replicaID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	actual, err := replicas.Get(context.TODO(), getClient("2.11"), replicaID).Extract()

	expected := &replicas.Replica{
		AvailabilityZone: "zone-1",
		ShareNetworkID:   "ca0163c8-3941-4420-8b01-41517e19e366",
		ShareServerID:    "5ccc1b0c-334a-4e46-81e6-b52e03223060",
		ShareID:          "65a34695-f9e5-4eea-b48d-a0b261d82943",
		ID:               replicaID,
		Status:           "available",
		State:            "active",
		CreatedAt:        time.Date(2023, time.May, 26, 12, 32, 56, 391337000, time.UTC),
		UpdatedAt:        time.Date(2023, time.May, 26, 12, 33, 28, 265716000, time.UTC),
	}

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	listOpts := &replicas.ListOpts{
		ShareID: "65a34695-f9e5-4eea-b48d-a0b261d82943",
	}
	allPages, err := replicas.List(getClient("2.11"), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := replicas.ExtractReplicas(allPages)
	th.AssertNoErr(t, err)

	expected := []replicas.Replica{
		{
			ID:      replicaID,
			ShareID: "65a34695-f9e5-4eea-b48d-a0b261d82943",
			Status:  "available",
			State:   "active",
		},
		{
			ID:      "4b70c2e2-eec7-4699-880d-4da9051ca162",
			ShareID: "65a34695-f9e5-4eea-b48d-a0b261d82943",
			Status:  "available",
			State:   "out_of_sync",
		},
		{
			ID:      "920bb037-bdd7-48a1-98f0-1aa1787ca3eb",
			ShareID: "65a34695-f9e5-4eea-b48d-a0b261d82943",
			Status:  "available",
			State:   "in_sync",
		},
	}

	th.AssertDeepEquals(t, expected, actual)
}

func TestListDetail(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListDetailResponse(t)

	listOpts := &replicas.ListOpts{
		ShareID: "65a34695-f9e5-4eea-b48d-a0b261d82943",
	}
	allPages, err := replicas.ListDetail(getClient("2.11"), listOpts).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	actual, err := replicas.ExtractReplicas(allPages)
	th.AssertNoErr(t, err)

	expected := []replicas.Replica{
		{
			AvailabilityZone: "zone-1",
			ShareNetworkID:   "ca0163c8-3941-4420-8b01-41517e19e366",
			ShareServerID:    "5ccc1b0c-334a-4e46-81e6-b52e03223060",
			ShareID:          "65a34695-f9e5-4eea-b48d-a0b261d82943",
			ID:               replicaID,
			Status:           "available",
			State:            "active",
			CreatedAt:        time.Date(2023, time.May, 26, 12, 32, 56, 391337000, time.UTC),
			UpdatedAt:        time.Date(2023, time.May, 26, 12, 33, 28, 265716000, time.UTC),
		},
		{
			AvailabilityZone: "zone-2",
			ShareNetworkID:   "ca0163c8-3941-4420-8b01-41517e19e366",
			ShareServerID:    "81aa586e-3a03-4f92-98bd-807d87a61c1a",
			ShareID:          "65a34695-f9e5-4eea-b48d-a0b261d82943",
			ID:               "4b70c2e2-eec7-4699-880d-4da9051ca162",
			Status:           "available",
			State:            "out_of_sync",
			CreatedAt:        time.Date(2023, time.May, 26, 11, 59, 38, 313089000, time.UTC),
			UpdatedAt:        time.Date(2023, time.May, 26, 12, 00, 04, 321081000, time.UTC),
		},
		{
			AvailabilityZone: "zone-1",
			ShareNetworkID:   "ca0163c8-3941-4420-8b01-41517e19e366",
			ShareServerID:    "b87ea601-7d4c-47f3-8956-6876b7a6b6db",
			ShareID:          "65a34695-f9e5-4eea-b48d-a0b261d82943",
			ID:               "920bb037-bdd7-48a1-98f0-1aa1787ca3eb",
			Status:           "available",
			State:            "in_sync",
			CreatedAt:        time.Date(2023, time.May, 26, 12, 32, 45, 751834000, time.UTC),
			UpdatedAt:        time.Date(2023, time.May, 26, 12, 36, 04, 110328000, time.UTC),
		},
	}

	th.AssertDeepEquals(t, expected, actual)
}

func TestListExportLocationsSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListExportLocationsResponse(t)

	actual, err := replicas.ListExportLocations(context.TODO(), getClient("2.47"), replicaID).Extract()

	expected := []replicas.ExportLocation{
		{
			ID:               "3fc02d3c-da47-42a2-88b8-2d48f8c276bd",
			Path:             "192.168.1.123:/var/lib/manila/mnt/share-3b9c33e8-b136-45c6-84a6-019c8db1d550",
			Preferred:        true,
			State:            "active",
			AvailabilityZone: "zone-1",
		},
		{
			ID:               "ae73e762-e8b9-4aad-aad3-23afb7cd6825",
			Path:             "192.168.1.124:/var/lib/manila/mnt/share-3b9c33e8-b136-45c6-84a6-019c8db1d550",
			Preferred:        false,
			State:            "active",
			AvailabilityZone: "zone-1",
		},
	}

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expected, actual)
}

func TestGetExportLocationSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetExportLocationResponse(t)

	s, err := replicas.GetExportLocation(context.TODO(), getClient("2.47"), replicaID, "ae73e762-e8b9-4aad-aad3-23afb7cd6825").Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, s, &replicas.ExportLocation{
		Path:             "192.168.1.124:/var/lib/manila/mnt/share-3b9c33e8-b136-45c6-84a6-019c8db1d550",
		ID:               "ae73e762-e8b9-4aad-aad3-23afb7cd6825",
		Preferred:        false,
		State:            "active",
		AvailabilityZone: "zone-1",
		CreatedAt:        time.Date(2023, time.May, 26, 12, 44, 33, 987960000, time.UTC),
		UpdatedAt:        time.Date(2023, time.May, 26, 12, 44, 33, 958363000, time.UTC),
	})
}

func TestResetStatusSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockResetStatusResponse(t)

	err := replicas.ResetStatus(context.TODO(), getClient("2.11"), replicaID, &replicas.ResetStatusOpts{Status: "available"}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestResetStateSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockResetStateResponse(t)

	err := replicas.ResetState(context.TODO(), getClient("2.11"), replicaID, &replicas.ResetStateOpts{State: "active"}).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestResyncSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockResyncResponse(t)

	err := replicas.Resync(context.TODO(), getClient("2.11"), replicaID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestPromoteSuccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockPromoteResponse(t)

	err := replicas.Promote(context.TODO(), getClient("2.11"), replicaID, &replicas.PromoteOpts{QuiesceWaitTime: 30}).ExtractErr()
	th.AssertNoErr(t, err)
}
