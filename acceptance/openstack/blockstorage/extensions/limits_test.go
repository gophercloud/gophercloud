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

	th.AssertEquals(t, limits.Absolute.MaxTotalVolumes, 40)
	th.AssertEquals(t, limits.Absolute.MaxTotalSnapshots, 40)
	th.AssertEquals(t, limits.Absolute.MaxTotalVolumeGigabytes, 1000)
	th.AssertEquals(t, limits.Absolute.MaxTotalBackups, 10)
	th.AssertEquals(t, limits.Absolute.MaxTotalBackupGigabytes, 1000)
	th.AssertEquals(t, limits.Absolute.TotalVolumesUsed, 0)
	th.AssertEquals(t, limits.Absolute.TotalGigabytesUsed, 0)
	th.AssertEquals(t, limits.Absolute.TotalSnapshotsUsed, 0)
	th.AssertEquals(t, limits.Absolute.TotalBackupsUsed, 0)
	th.AssertEquals(t, limits.Absolute.TotalBackupGigabytesUsed, 0)
}
