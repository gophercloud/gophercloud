package datastores

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
	"github.com/rackspace/gophercloud/testhelper/fixture"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, "/datastores", "GET", "", listDSResp, 200)

	pages := 0

	err := List(fake.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := ExtractDatastores(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, []Datastore{exampleDatastore}, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, "/datastores/{dsID}", "GET", "", getDSResp, 200)

	ds, err := Get(fake.ServiceClient(), "{dsID}").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &exampleDatastore, ds)
}

func TestListVersions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, "/datastores/{dsID}/versions", "GET", "", listVersionsResp, 200)

	pages := 0

	err := ListVersions(fake.ServiceClient(), "{dsID}").EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := ExtractVersions(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, exampleVersions, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestGetVersion(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixture.SetupHandler(t, "/datastores/{dsID}/versions/{versionID}", "GET", "", getVersionResp, 200)

	ds, err := GetVersion(fake.ServiceClient(), "{dsID}", "{versionID}").Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, &exampleVersion1, ds)
}
