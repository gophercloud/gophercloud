package testing

import (
	"context"
	"testing"

	db "github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreate(t, fakeServer)

	opts := users.BatchCreateOpts{
		{
			Databases: db.BatchCreateOpts{
				db.CreateOpts{Name: "databaseA"},
			},
			Name:     "dbuser3",
			Password: "secretsecret",
		},
		{
			Databases: db.BatchCreateOpts{
				{Name: "databaseB"},
				{Name: "databaseC"},
			},
			Name:     "dbuser4",
			Password: "secretsecret",
		},
	}

	res := users.Create(context.TODO(), client.ServiceClient(fakeServer), instanceID, opts)
	th.AssertNoErr(t, res.Err)
}

func TestUserList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleList(t, fakeServer)

	expectedUsers := []users.User{
		{
			Databases: []db.Database{
				{Name: "databaseA"},
			},
			Name: "dbuser3",
		},
		{
			Databases: []db.Database{
				{Name: "databaseB"},
				{Name: "databaseC"},
			},
			Name: "dbuser4",
		},
	}

	pages := 0
	err := users.List(client.ServiceClient(fakeServer), instanceID).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := users.ExtractUsers(page)
		if err != nil {
			return false, err
		}

		th.CheckDeepEquals(t, expectedUsers, actual)
		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, pages)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDelete(t, fakeServer)

	res := users.Delete(context.TODO(), client.ServiceClient(fakeServer), instanceID, "{userName}")
	th.AssertNoErr(t, res.Err)
}
