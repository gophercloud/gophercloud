package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestListUsers(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListUsersSuccessfully(t)

	count := 0
	err := users.List(client.ServiceClient(), nil).EachPage(func(page pagination.Page) (bool, error) {
		count++

		actual, err := users.ExtractUsers(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedUsersSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, count, 1)
}

func TestListUsersAllPages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListUsersSuccessfully(t)

	allPages, err := users.List(client.ServiceClient(), nil).AllPages()
	th.AssertNoErr(t, err)
	actual, err := users.ExtractUsers(allPages)
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedUsersSlice, actual)
}

func TestGetUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetUserSuccessfully(t)

	actual, err := users.Get(client.ServiceClient(), "9fe1d3", nil).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondUser, *actual)
}

func TestCreateUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateUserSuccessfully(t)

	iTrue := true
	createOpts := users.CreateOpts{
		Name:             "glance",
		DomainID:         "default",
		Enabled:          &iTrue,
		Password:         "secretsecret",
		DefaultProjectID: "263fd9",
	}

	actual, err := users.Create(client.ServiceClient(), createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, SecondUser, *actual)
}
