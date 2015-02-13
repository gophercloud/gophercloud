package users

import (
	"testing"

	db "github.com/rackspace/gophercloud/openstack/db/v1/databases"
	os "github.com/rackspace/gophercloud/openstack/db/v1/users"
	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	fake "github.com/rackspace/gophercloud/testhelper/client"
)

const userName = "{userName}"

func TestChangeUserPassword(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleChangePasswordSuccessfully(t, instanceID)

	opts := os.BatchCreateOpts{
		os.CreateOpts{Name: "dbuser1", Password: "newpassword"},
		os.CreateOpts{Name: "dbuser2", Password: "anotherpassword"},
	}

	err := ChangePassword(fake.ServiceClient(), instanceID, opts).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestUpdateUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleUpdateSuccessfully(t, instanceID, userName)

	opts := os.CreateOpts{
		Name:     "new_username",
		Password: "new_password",
	}

	err := Update(fake.ServiceClient(), instanceID, userName, opts).ExtractErr()

	th.AssertNoErr(t, err)
}

func TestGetUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGetSuccessfully(t, instanceID, userName)

	user, err := Get(fake.ServiceClient(), instanceID, userName).Extract()

	th.AssertNoErr(t, err)

	expected := &User{
		Name: "exampleuser",
		Host: "foo",
		Databases: []db.Database{
			db.Database{Name: "databaseA"},
			db.Database{Name: "databaseB"},
		},
	}

	th.AssertDeepEquals(t, expected, user)
}

func TestUserAccessList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleListUserAccessSuccessfully(t, instanceID, userName)

	expectedDBs := []db.Database{
		db.Database{Name: "databaseE"},
	}

	pages := 0
	err := ListAccess(fake.ServiceClient(), instanceID, userName).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := ExtractDBs(page)
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

func TestGrantAccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleGrantUserAccessSuccessfully(t, instanceID, userName)

	opts := db.BatchCreateOpts{
		db.CreateOpts{Name: "databaseE"},
	}

	err := GrantAccess(fake.ServiceClient(), instanceID, userName, opts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestRevokeAccess(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	HandleRevokeUserAccessSuccessfully(t, instanceID, userName, "{dbName}")

	err := RevokeAccess(fake.ServiceClient(), instanceID, userName, "{dbName}").ExtractErr()
	th.AssertNoErr(t, err)
}
