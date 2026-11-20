package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/datastores"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
	"github.com/gophercloud/gophercloud/v2/testhelper/fixture"
)

func TestCreateVersion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, "/mgmt/datastore-versions", "POST", createVersionReq, CreateVersionResp, 202)

	opts := datastores.CreateVersionOpts{
		Name:             "mysql-5.7",
		DatastoreName:    "mysql",
		DatastoreManager: "mysql",
		ImageTags:        []string{"trove", "mysql"},
		Active:           gophercloud.Enabled,
		Default:          gophercloud.Disabled,
		Packages:         []string{"mysql-server"},
		Version:          "5.7.30",
	}

	version, err := datastores.CreateVersion(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExampleCreatedVersion, version)
}

func TestListAllVersions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, "/mgmt/datastore-versions", "GET", "", ListAllVersionsResp, 200)

	pages := 0
	err := datastores.ListAllVersions(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := datastores.ExtractVersions(page)
		if err != nil {
			return false, err
		}

		th.AssertDeepEquals(t, []datastores.Version{ExampleCreatedVersion}, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGetVersionByID(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, "/mgmt/datastore-versions/{versionID}", "GET", "", CreateVersionResp, 200)

	version, err := datastores.GetVersionByID(context.TODO(), client.ServiceClient(fakeServer), "{versionID}").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExampleCreatedVersion, version)
}

func TestUpdateVersion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, "/mgmt/datastore-versions/{versionID}", "PATCH", updateVersionReq, UpdateVersionResp, 202)

	opts := datastores.UpdateVersionOpts{
		Name:   "mysql-5.7-updated",
		Active: gophercloud.Enabled,
	}

	version, err := datastores.UpdateVersion(context.TODO(), client.ServiceClient(fakeServer), "{versionID}", opts).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExampleUpdatedVersion, version)
}

func TestDeleteVersion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, "/mgmt/datastore-versions/{versionID}", "DELETE", "", "", 202)

	err := datastores.DeleteVersion(context.TODO(), client.ServiceClient(fakeServer), "{versionID}").ExtractErr()
	th.AssertNoErr(t, err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, "/datastores", "GET", "", ListDSResp, 200)

	pages := 0

	err := datastores.List(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := datastores.ExtractDatastores(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, []datastores.Datastore{ExampleDatastore}, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, "/datastores/{dsID}", "GET", "", GetDSResp, 200)

	ds, err := datastores.Get(context.TODO(), client.ServiceClient(fakeServer), "{dsID}").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExampleDatastore, ds)
}

func TestListVersions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, "/datastores/{dsID}/versions", "GET", "", ListVersionsResp, 200)

	pages := 0

	err := datastores.ListVersions(client.ServiceClient(fakeServer), "{dsID}").EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := datastores.ExtractVersions(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, ExampleVersions, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGetVersion(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	fixture.SetupHandler(t, fakeServer, "/datastores/{dsID}/versions/{versionID}", "GET", "", GetVersionResp, 200)

	ds, err := datastores.GetVersion(context.TODO(), client.ServiceClient(fakeServer), "{dsID}", "{versionID}").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &ExampleVersion1, ds)
}
