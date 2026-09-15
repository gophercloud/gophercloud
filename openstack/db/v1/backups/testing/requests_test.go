package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/backups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
	"github.com/gophercloud/gophercloud/v2/testhelper/fixture"
)

const (
	rootURL     = "/backups"
	backupID    = "a9832168-7541-4536-b8d9-a8a9b79cf1b4"
	resourceURL = rootURL + "/" + backupID
	instanceID  = "44b277eb-39be-4921-be31-3d61b43651d7"
	instanceURL = "/instances/" + instanceID + "/backups"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, rootURL, "GET", "", ListBackupsJSON, 200)

	pages := 0
	err := backups.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := backups.ExtractBackups(page)
		if err != nil {
			return false, err
		}

		th.AssertDeepEquals(t, []backups.Backup{ExampleBackup}, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestListByInstance(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, instanceURL, "GET", "", ListBackupsJSON, 200)

	pages := 0
	err := backups.ListByInstance(client.ServiceClient(fakeServer), instanceID).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := backups.ExtractBackups(page)
		if err != nil {
			return false, err
		}

		th.AssertDeepEquals(t, []backups.Backup{ExampleBackup}, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, rootURL, "POST", CreateBackupReq, CreateBackupJSON, 202)

	incremental := 0
	opts := backups.CreateOpts{
		Name:          "snapshot",
		InstanceID:    instanceID,
		Description:   "My Backup",
		Incremental:   &incremental,
		StorageDriver: "swift",
	}

	backup, err := backups.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExampleBackup, backup)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, resourceURL, "GET", "", GetBackupJSON, 200)

	backup, err := backups.Get(context.TODO(), client.ServiceClient(fakeServer), backupID).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExampleBackup, backup)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, resourceURL, "DELETE", "", "", 202)

	err := backups.Delete(context.TODO(), client.ServiceClient(fakeServer), backupID).ExtractErr()
	th.AssertNoErr(t, err)
}
