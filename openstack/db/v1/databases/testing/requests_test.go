package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreate(t, fakeServer)

	opts := databases.BatchCreateOpts{
		databases.CreateOpts{Name: "testingdb", CharSet: "utf8", Collate: "utf8_general_ci"},
		databases.CreateOpts{Name: "sampledb"},
	}

	res := databases.Create(context.TODO(), client.ServiceClient(fakeServer), instanceID, opts)
	th.AssertNoErr(t, res.Err)
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleList(t, fakeServer)

	expectedDBs := []databases.Database{
		{Name: "anotherexampledb"},
		{Name: "exampledb"},
		{Name: "nextround"},
		{Name: "sampledb"},
		{Name: "testingdb"},
	}

	pages := 0
	err := databases.List(client.ServiceClient(fakeServer), instanceID).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := databases.ExtractDBs(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, expectedDBs, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDelete(t, fakeServer)

	err := databases.Delete(context.TODO(), client.ServiceClient(fakeServer), instanceID, "{dbName}").ExtractErr()
	th.AssertNoErr(t, err)
}
