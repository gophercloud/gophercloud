package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/db/v1/backup"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreate(t)

	opts := backup.CreateOpts{
		Name:        "testingbackup",
		Instance:    instanceID,
		Incremental: 0,
	}

	backup, err := backup.Create(fake.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectBackup, backup)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleList(t)

	expectedBackups := []backup.Backup{expectBackup}

	pages := 0
	err := backup.List(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := backup.ExtractBackups(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, expectedBackups, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGet(t)

	backup, err := backup.Get(fake.ServiceClient(), backupID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectBackup, backup)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDelete(t)

	err := backup.Delete(fake.ServiceClient(), backupID).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestListByInstance(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListByInstance(t)

	expectedBackups := []backup.Backup{expectBackup}

	pages := 0
	err := backup.ListByInstance(fake.ServiceClient(), instanceID).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := backup.ExtractBackups(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, expectedBackups, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}
