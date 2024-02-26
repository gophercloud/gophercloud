package testing

import (
	"context"
	"testing"

	db "github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreate(t)

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

	res := users.Create(context.TODO(), fake.ServiceClient(), instanceID, opts)
	th.AssertNoErr(t, res.Err)
}

func TestUserList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleList(t)

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
	err := users.List(fake.ServiceClient(), instanceID).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDelete(t)

	res := users.Delete(context.TODO(), fake.ServiceClient(), instanceID, "{userName}")
	th.AssertNoErr(t, res.Err)
}
