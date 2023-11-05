package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/db/v1/backups"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreate(t)

	opts := backups.CreateOpts{
		Name:        "snapshot",
		Instance:    "44b277eb-39be-4921-be31-3d61b43651d7",
		Description: "My Backup",
		Incremental: 0,
	}

	backup, err := backups.Create(fake.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedBackup, backup)
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleList(t)

	pages := 0
	err := backups.List(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := backups.ExtractBackups(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, []backups.Backup{expectedBackup}, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGet(t)

	backup, err := backups.Get(fake.ServiceClient(), backupId).Extract()

	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &expectedGetBackup, backup)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDelete(t)

	res := backups.Delete(fake.ServiceClient(), backupId)
	th.AssertNoErr(t, res.Err)
}
