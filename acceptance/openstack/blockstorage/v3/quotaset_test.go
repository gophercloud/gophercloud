// +build acceptance quotasets

package v3

import (
	"fmt"
	"os"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/extensions/quotasets"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestQuotasetGet(t *testing.T) {
	client, projectID := getClientAndProject(t)

	quotaSet, err := quotasets.Get(client, projectID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaSet)
}

func TestQuotasetGetDefaults(t *testing.T) {
	client, projectID := getClientAndProject(t)

	quotaSet, err := quotasets.GetDefaults(client, projectID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, quotaSet)
}

func TestQuotasetGetDetailed(t *testing.T) {
	client, projectID := getClientAndProject(t)

	detailedQuotaSet, err := quotasets.GetDetail(client, projectID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, detailedQuotaSet)
}

var UpdateQuotaOpts = quotasets.UpdateOpts{
	Volumes:            gophercloud.IntToPointer(100),
	Snapshots:          gophercloud.IntToPointer(200),
	Gigabytes:          gophercloud.IntToPointer(300),
	PerVolumeGigabytes: gophercloud.IntToPointer(50),
	Backups:            gophercloud.IntToPointer(2),
	BackupGigabytes:    gophercloud.IntToPointer(300),
}

var UpdatedQuotas = quotasets.QuotaSet{
	Volumes:            100,
	Snapshots:          200,
	Gigabytes:          300,
	PerVolumeGigabytes: 50,
	Backups:            2,
	BackupGigabytes:    300,
}

func TestQuotasetUpdate(t *testing.T) {
	client, projectID := getClientAndProject(t)

	// save original quotas
	orig, err := quotasets.Get(client, projectID).Extract()
	th.AssertNoErr(t, err)

	defer func() {
		restore := quotasets.UpdateOpts{}
		FillUpdateOptsFromQuotaSet(*orig, &restore)

		_, err = quotasets.Update(client, projectID, restore).Extract()
		th.AssertNoErr(t, err)
	}()

	// test Update
	resultQuotas, err := quotasets.Update(client, projectID, UpdateQuotaOpts).Extract()
	th.AssertNoErr(t, err)
	th.AssertEquals(t, UpdatedQuotas, *resultQuotas)

	// We dont know the default quotas, so just check if the quotas are not the
	// same as before
	newQuotas, err := quotasets.Get(client, projectID).Extract()
	th.AssertNoErr(t, err)
	if newQuotas.Volumes != resultQuotas.Volumes {
		t.Fatalf(
			fmt.Sprintf("Failed to update quotas\n\texpected: %d\t actual: %d",
				resultQuotas.Volumes,
				newQuotas.Volumes,
			),
		)
	}
}

func TestQuotasetDelete(t *testing.T) {
	client, projectID := getClientAndProject(t)

	// save original quotas
	orig, err := quotasets.Get(client, projectID).Extract()
	th.AssertNoErr(t, err)

	defer func() {
		restore := quotasets.UpdateOpts{}
		FillUpdateOptsFromQuotaSet(*orig, &restore)

		_, err = quotasets.Update(client, projectID, restore).Extract()
		th.AssertNoErr(t, err)
	}()

	// Obtain environment default quotaset values to validate deletion.
	defaultQuotaSet, err := quotasets.GetDefaults(client, projectID).Extract()
	th.AssertNoErr(t, err)

	// Test Delete
	_, err = quotasets.Delete(client, projectID).Extract()
	th.AssertNoErr(t, err)

	newQuotas, err := quotasets.Get(client, projectID).Extract()
	th.AssertNoErr(t, err)

	if newQuotas.Volumes != defaultQuotaSet.Volumes {
		t.Fatalf("Failed to delete quotas!")
	}
}

// getClientAndProject reduces boilerplate by returning a new blockstorage v3
// ServiceClient and a project ID obtained from the OS_PROJECT_NAME envvar.
func getClientAndProject(t *testing.T) (*gophercloud.ServiceClient, string) {
	client, err := clients.NewBlockStorageV3Client()
	th.AssertNoErr(t, err)

	projectID := os.Getenv("OS_PROJECT_NAME")
	th.AssertNoErr(t, err)
	return client, projectID
}
func FillUpdateOptsFromQuotaSet(src quotasets.QuotaSet, dest *quotasets.UpdateOpts) {
	dest.Volumes = &src.Volumes
	dest.Snapshots = &src.Snapshots
	dest.Gigabytes = &src.Gigabytes
	dest.PerVolumeGigabytes = &src.PerVolumeGigabytes
	dest.Backups = &src.Backups
	dest.BackupGigabytes = &src.BackupGigabytes
}
