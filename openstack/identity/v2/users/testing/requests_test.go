package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v2/users"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	MockListUserResponse(t, fakeServer)

	count := 0

	err := users.List(client.ServiceClient(fakeServer)).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := users.ExtractUsers(page)
		th.AssertNoErr(t, err)

		expected := []users.User{
			{
				ID:       "u1000",
				Name:     "John Smith",
				Username: "jqsmith",
				Email:    "john.smith@example.org",
				Enabled:  true,
				TenantID: "12345",
			},
			{
				ID:       "u1001",
				Name:     "Jane Smith",
				Username: "jqsmith",
				Email:    "jane.smith@example.org",
				Enabled:  true,
				TenantID: "12345",
			},
		}
		th.CheckDeepEquals(t, expected, actual)
		return true, nil
	})
	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestCreateUser(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockCreateUserResponse(t, fakeServer)

	opts := users.CreateOpts{
		Name:     "new_user",
		TenantID: "12345",
		Enabled:  gophercloud.Disabled,
		Email:    "new_user@foo.com",
	}

	user, err := users.Create(context.TODO(), client.ServiceClient(fakeServer), opts).Extract()

	th.AssertNoErr(t, err)

	expected := &users.User{
		Name:     "new_user",
		ID:       "c39e3de9be2d4c779f1dfd6abacc176d",
		Email:    "new_user@foo.com",
		Enabled:  false,
		TenantID: "12345",
	}

	th.AssertDeepEquals(t, expected, user)
}

func TestGetUser(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockGetUserResponse(t, fakeServer)

	user, err := users.Get(context.TODO(), client.ServiceClient(fakeServer), "new_user").Extract()
	th.AssertNoErr(t, err)

	expected := &users.User{
		Name:     "new_user",
		ID:       "c39e3de9be2d4c779f1dfd6abacc176d",
		Email:    "new_user@foo.com",
		Enabled:  false,
		TenantID: "12345",
	}

	th.AssertDeepEquals(t, expected, user)
}

func TestUpdateUser(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockUpdateUserResponse(t, fakeServer)

	id := "c39e3de9be2d4c779f1dfd6abacc176d"
	opts := users.UpdateOpts{
		Name:    "new_name",
		Enabled: gophercloud.Enabled,
		Email:   "new_email@foo.com",
	}

	user, err := users.Update(context.TODO(), client.ServiceClient(fakeServer), id, opts).Extract()

	th.AssertNoErr(t, err)

	expected := &users.User{
		Name:     "new_name",
		ID:       id,
		Email:    "new_email@foo.com",
		Enabled:  true,
		TenantID: "12345",
	}

	th.AssertDeepEquals(t, expected, user)
}

func TestDeleteUser(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockDeleteUserResponse(t, fakeServer)

	res := users.Delete(context.TODO(), client.ServiceClient(fakeServer), "c39e3de9be2d4c779f1dfd6abacc176d")
	th.AssertNoErr(t, res.Err)
}

func TestListingUserRoles(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	mockListRolesResponse(t, fakeServer)

	tenantID := "1d8b6120dcc640fda4fc9194ffc80273"
	userID := "c39e3de9be2d4c779f1dfd6abacc176d"

	err := users.ListRoles(client.ServiceClient(fakeServer), tenantID, userID).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		actual, err := users.ExtractRoles(page)
		th.AssertNoErr(t, err)

		expected := []users.Role{
			{ID: "9fe2ff9ee4384b1894a90878d3e92bab", Name: "foo_role"},
			{ID: "1ea3d56793574b668e85960fbf651e13", Name: "admin"},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
}
