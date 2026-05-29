//go:build acceptance || blockstorage || backups

package v3

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/backups"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestBackupsCRUD(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	backup, err := CreateBackup(t, blockClient, volume.ID)
	th.AssertNoErr(t, err)
	defer DeleteBackup(t, blockClient, backup.ID)

	allPages, err := backups.List(blockClient, nil).AllPages(context.TODO())
	th.AssertNoErr(t, err)

	allBackups, err := backups.ExtractBackups(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, v := range allBackups {
		if backup.Name == v.Name {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestBackupsResetStatus(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	backup, err := CreateBackup(t, blockClient, volume.ID)
	th.AssertNoErr(t, err)
	defer DeleteBackup(t, blockClient, backup.ID)

	err = ResetBackupStatus(t, blockClient, backup, "error")
	th.AssertNoErr(t, err)

	err = ResetBackupStatus(t, blockClient, backup, "available")
	th.AssertNoErr(t, err)
}

func TestBackupsForceDelete(t *testing.T) {
	blockClient, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	volume, err := CreateVolume(t, blockClient)
	th.AssertNoErr(t, err)
	defer DeleteVolume(t, blockClient, volume)

	backup, err := CreateBackup(t, blockClient, volume.ID)
	th.AssertNoErr(t, err)
	defer DeleteBackup(t, blockClient, backup.ID)

	err = WaitForBackupStatus(blockClient, backup.ID, "available")
	th.AssertNoErr(t, err)

	err = backups.ForceDelete(context.TODO(), blockClient, backup.ID).ExtractErr()
	th.AssertNoErr(t, err)

	err = WaitForBackupStatus(blockClient, backup.ID, "deleted")
	th.AssertNoErr(t, err)
}
