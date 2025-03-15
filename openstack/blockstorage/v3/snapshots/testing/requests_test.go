package testing

import (
	"context"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	count := 0

	err := snapshots.List(client.ServiceClient(), &snapshots.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := snapshots.ExtractSnapshots(page)
		if err != nil {
			t.Errorf("Failed to extract snapshots: %v", err)
			return false, err
		}

		expected := []snapshots.Snapshot{
			{
				ID:          "289da7f8-6440-407c-9fb4-7db01ec49164",
				Name:        "snapshot-001",
				VolumeID:    "521752a6-acf6-4b2d-bc7a-119f9148cd8c",
				Status:      "available",
				Size:        30,
				CreatedAt:   time.Date(2017, 5, 30, 3, 35, 3, 0, time.UTC),
				Description: "Daily Backup",
			},
			{
				ID:          "96c3bda7-c82a-4f50-be73-ca7621794835",
				Name:        "snapshot-002",
				VolumeID:    "76b8950a-8594-4e5b-8dce-0dfa9c696358",
				Status:      "available",
				Size:        25,
				CreatedAt:   time.Date(2017, 5, 30, 3, 35, 3, 0, time.UTC),
				Description: "Weekly Backup",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestDetailList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListDetailsResponse(t)

	count := 0

	err := snapshots.ListDetail(client.ServiceClient(), &snapshots.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := snapshots.ExtractSnapshots(page)
		if err != nil {
			t.Errorf("Failed to extract snapshots: %v", err)
			return false, err
		}
		expected := []snapshots.Snapshot{
			{
				ID:              "289da7f8-6440-407c-9fb4-7db01ec49164",
				Name:            "snapshot-001",
				VolumeID:        "521752a6-acf6-4b2d-bc7a-119f9148cd8c",
				Status:          "available",
				Size:            30,
				CreatedAt:       time.Date(2017, 5, 30, 3, 35, 3, 0, time.UTC),
				Description:     "Daily Backup",
				Progress:        "100%",
				ProjectID:       "84b8950a-8594-4e5b-8dce-0dfa9c696357",
				GroupSnapshotID: "",
				UserID:          "075da7f8-6440-407c-9fb4-7db01ec49531",
				ConsumesQuota:   true,
			},
			{
				ID:              "96c3bda7-c82a-4f50-be73-ca7621794835",
				Name:            "snapshot-002",
				VolumeID:        "76b8950a-8594-4e5b-8dce-0dfa9c696358",
				Status:          "available",
				Size:            25,
				CreatedAt:       time.Date(2017, 5, 30, 3, 35, 3, 0, time.UTC),
				Description:     "Weekly Backup",
				Progress:        "50%",
				ProjectID:       "84b8950a-8594-4e5b-8dce-0dfa9c696357",
				GroupSnapshotID: "865da7f8-6440-407c-9fb4-7db01ec40876",
				UserID:          "075da7f8-6440-407c-9fb4-7db01ec49531",
				ConsumesQuota:   false,
			},
		}
		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	v, err := snapshots.Get(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, v.Name, "snapshot-001")
	th.AssertEquals(t, v.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := snapshots.CreateOpts{VolumeID: "1234", Name: "snapshot-001"}
	n, err := snapshots.Create(context.TODO(), client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.VolumeID, "1234")
	th.AssertEquals(t, n.Name, "snapshot-001")
	th.AssertEquals(t, n.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestUpdateMetadata(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateMetadataResponse(t)

	expected := map[string]any{"key": "v1"}

	options := &snapshots.UpdateMetadataOpts{
		Metadata: map[string]any{
			"key": "v1",
		},
	}

	actual, err := snapshots.UpdateMetadata(context.TODO(), client.ServiceClient(), "123", options).ExtractMetadata()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := snapshots.Delete(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(t, res.Err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateResponse(t)

	var name = "snapshot-002"
	var description = "Daily backup 002"
	options := snapshots.UpdateOpts{Name: &name, Description: &description}
	v, err := snapshots.Update(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", options).Extract()
	th.AssertNoErr(t, err)
	th.CheckEquals(t, "snapshot-002", v.Name)
	th.CheckEquals(t, "Daily backup 002", v.Description)
}

func TestResetStatus(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockResetStatusResponse(t)

	opts := &snapshots.ResetStatusOpts{
		Status: "error",
	}
	res := snapshots.ResetStatus(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", opts)
	th.AssertNoErr(t, res.Err)
}

func TestUpdateStatus(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateStatusResponse(t)

	opts := &snapshots.UpdateStatusOpts{
		Status:   "error",
		Progress: "80%",
	}
	res := snapshots.UpdateStatus(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", opts)
	th.AssertNoErr(t, res.Err)
}

func TestForceDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockForceDeleteResponse(t)

	res := snapshots.ForceDelete(context.TODO(), client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(t, res.Err)
}
