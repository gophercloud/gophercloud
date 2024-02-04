package testing

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/snapshotextensions"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListAllWithExtensions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	type SnapshotWithExt struct {
		snapshots.Snapshot
		snapshotextensions.SnapshotExtension
	}

	allPages, err := snapshots.List(client.ServiceClient(), &snapshots.ListOpts{}).AllPages()
	th.AssertNoErr(t, err)

	var actual []SnapshotWithExt
	err = snapshots.ExtractSnapshotsInto(allPages, &actual)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 2, len(actual))
	th.AssertEquals(t, "89afd400-b646-4bbc-b12b-c0a4d63e5bd3", actual[0].ProjectId)
	th.AssertEquals(t, "89afd400-b646-4bbc-b12b-c0a4d63e5bd3", actual[1].ProjectId)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	count := 0

	snapshots.List(client.ServiceClient(), &snapshots.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
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

	if count != 1 {
		t.Errorf("Expected 1 page, got %d", count)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	v, err := snapshots.Get(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, v.Name, "snapshot-001")
	th.AssertEquals(t, v.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestGetWithExtensions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	var s struct {
		snapshots.Snapshot
		snapshotextensions.SnapshotExtension
	}
	err := snapshots.Get(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").ExtractInto(&s)
	th.AssertNoErr(t, err)
	th.AssertEquals(t, "89afd400-b646-4bbc-b12b-c0a4d63e5bd3", s.ProjectId)

	err = snapshots.Get(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").ExtractInto(s)
	if err == nil {
		t.Errorf("Expected error when providing non-pointer struct")
	}
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := snapshots.CreateOpts{VolumeID: "1234", Name: "snapshot-001"}
	n, err := snapshots.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.VolumeID, "1234")
	th.AssertEquals(t, n.Name, "snapshot-001")
	th.AssertEquals(t, n.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestUpdateMetadata(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateMetadataResponse(t)

	expected := map[string]interface{}{"key": "v1"}

	options := &snapshots.UpdateMetadataOpts{
		Metadata: map[string]interface{}{
			"key": "v1",
		},
	}

	actual, err := snapshots.UpdateMetadata(client.ServiceClient(), "123", options).ExtractMetadata()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, actual, expected)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := snapshots.Delete(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(t, res.Err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateResponse(t)

	var name = "snapshot-002"
	var description = "Daily backup 002"
	options := snapshots.UpdateOpts{Name: &name, Description: &description}
	v, err := snapshots.Update(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", options).Extract()
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
	res := snapshots.ResetStatus(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", opts)
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
	res := snapshots.UpdateStatus(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", opts)
	th.AssertNoErr(t, res.Err)
}

func TestForceDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockForceDeleteResponse(t)

	res := snapshots.ForceDelete(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(t, res.Err)
}
