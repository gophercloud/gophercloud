package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/datastores"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
	"github.com/gophercloud/gophercloud/v2/testhelper/fixture"
)

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
