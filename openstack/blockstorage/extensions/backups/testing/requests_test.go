package testing

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/backups"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	count := 0

	backups.List(client.ServiceClient(), &backups.ListOpts{}).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := backups.ExtractBackups(page)
		if err != nil {
			t.Errorf("Failed to extract backups: %v", err)
			return false, err
		}

		expected := []backups.Backup{
			{
				ID:          "289da7f8-6440-407c-9fb4-7db01ec49164",
				Name:        "backup-001",
				VolumeID:    "521752a6-acf6-4b2d-bc7a-119f9148cd8c",
				Status:      "available",
				Size:        30,
				CreatedAt:   time.Date(2017, 5, 30, 3, 35, 3, 0, time.UTC),
				Description: "Daily Backup",
			},
			{
				ID:          "96c3bda7-c82a-4f50-be73-ca7621794835",
				Name:        "backup-002",
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

	v, err := backups.Get(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, v.Name, "backup-001")
	th.AssertEquals(t, v.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := backups.CreateOpts{VolumeID: "1234", Name: "backup-001"}
	n, err := backups.Create(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.VolumeID, "1234")
	th.AssertEquals(t, n.Name, "backup-001")
	th.AssertEquals(t, n.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestRestore(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockRestoreResponse(t)

	options := backups.RestoreOpts{VolumeID: "1234", Name: "vol-001"}
	n, err := backups.RestoreFromBackup(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22", options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.VolumeID, "1234")
	th.AssertEquals(t, n.VolumeName, "vol-001")
	th.AssertEquals(t, n.BackupID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := backups.Delete(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22")
	th.AssertNoErr(t, res.Err)
}

func TestExport(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockExportResponse(t)

	n, err := backups.Export(client.ServiceClient(), "d32019d3-bc6e-4319-9c1d-6722fc136a22").Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.BackupService, "cinder.backup.drivers.swift.SwiftBackupDriver")
	th.AssertDeepEquals(t, n.BackupURL, backupURL)

	tmp := backups.ImportBackup{}
	err = json.Unmarshal(backupURL, &tmp)
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, tmp, backupImport)
}

func TestImport(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockImportResponse(t)

	options := backups.ImportOpts{
		BackupService: "cinder.backup.drivers.swift.SwiftBackupDriver",
		BackupURL:     backupURL,
	}
	n, err := backups.Import(client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertEquals(t, n.ID, "d32019d3-bc6e-4319-9c1d-6722fc136a22")
}
