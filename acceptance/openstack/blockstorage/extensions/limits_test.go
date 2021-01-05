package extensions

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/limits"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestLimits(t *testing.T) {
	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	limits, err := limits.Get(client).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, limits)

	th.AssertIntGreaterOrEqual(t, limits.Absolute.MaxTotalVolumes, 0)
	th.AssertIntGreaterOrEqual(t, limits.Absolute.MaxTotalSnapshots, 0)
	th.AssertIntGreaterOrEqual(t, limits.Absolute.MaxTotalVolumeGigabytes, 0)
	th.AssertIntGreaterOrEqual(t, limits.Absolute.MaxTotalBackups, 0)
	th.AssertIntGreaterOrEqual(t, limits.Absolute.MaxTotalBackupGigabytes, 0)
	th.AssertIntGreaterOrEqual(t, limits.Absolute.TotalVolumesUsed, 0)
	th.AssertIntLesserOrEqual(t, limits.Absolute.TotalVolumesUsed, limits.Absolute.MaxTotalVolumes)
	th.AssertIntGreaterOrEqual(t, limits.Absolute.TotalGigabytesUsed, 0)
	th.AssertIntLesserOrEqual(t, limits.Absolute.TotalGigabytesUsed, limits.Absolute.MaxTotalVolumeGigabytes)
	th.AssertIntGreaterOrEqual(t, limits.Absolute.TotalSnapshotsUsed, 0)
	th.AssertIntLesserOrEqual(t, limits.Absolute.TotalSnapshotsUsed, limits.Absolute.MaxTotalSnapshots)
	th.AssertIntGreaterOrEqual(t, limits.Absolute.TotalBackupsUsed, 0)
	th.AssertIntLesserOrEqual(t, limits.Absolute.TotalBackupsUsed, limits.Absolute.MaxTotalBackups)
	th.AssertIntGreaterOrEqual(t, limits.Absolute.TotalBackupGigabytesUsed, 0)
	th.AssertIntLesserOrEqual(t, limits.Absolute.TotalBackupGigabytesUsed, limits.Absolute.MaxTotalBackupGigabytes)
}
