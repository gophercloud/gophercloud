package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/backupstrategies"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
	"github.com/gophercloud/gophercloud/v2/testhelper/fixture"
)

const (
	rootURL    = "/backup_strategies"
	instanceID = "0602db72-c63d-11ea-b87c-00224d6b7bc1"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, rootURL, "GET", "", ListBackupStrategiesJSON, 200)

	pages := 0
	err := backupstrategies.List(client.ServiceClient(fakeServer), nil).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := backupstrategies.ExtractBackupStrategies(page)
		if err != nil {
			return false, err
		}

		th.AssertDeepEquals(t, []backupstrategies.BackupStrategy{ExampleBackupStrategy}, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, rootURL, "POST", CreateBackupStrategyReq, CreateBackupStrategyJSON, 200)

	opts := backupstrategies.CreateOpts{
		InstanceID:     instanceID,
		SwiftContainer: "my_trove_backups",
	}

	strategy, err := backupstrategies.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExampleBackupStrategy, strategy)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, rootURL, "DELETE", "", "", 202)

	err := backupstrategies.Delete(context.TODO(), client.ServiceClient(fakeServer), nil).ExtractErr()
	th.AssertNoErr(t, err)
}
