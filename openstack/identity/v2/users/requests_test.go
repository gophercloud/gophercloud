package users

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	count := 0

	err := List(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractUsers(page)
		if err != nil {
			t.Errorf("Failed to extract users: %v", err)
			return false, err
		}

		expected := []User{
			User{
				ID:       "u1000",
				Name:     "John Smith",
				Username: "jqsmith",
				Email:    "john.smith@example.org",
				Enabled:  true,
				TenantID: "12345",
			},
			User{
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
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateUser(t)

	opts := CreateOpts{
		Name:     "new_user",
		TenantID: "12345",
		Enabled:  Disabled,
		Email:    "new_user@foo.com",
	}

	user, err := Create(client.ServiceClient(), opts).Extract()

	th.AssertNoErr(t, err)

	expected := &User{
		Name:     "new_user",
		ID:       "c39e3de9be2d4c779f1dfd6abacc176d",
		Email:    "new_user@foo.com",
		Enabled:  false,
		TenantID: "12345",
	}

	th.AssertDeepEquals(t, expected, user)
}

func TestGetUser(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetUser(t)

	user, err := Get(client.ServiceClient(), "new_user").Extract()
	th.AssertNoErr(t, err)

	expected := &User{
		Name:     "new_user",
		ID:       "c39e3de9be2d4c779f1dfd6abacc176d",
		Email:    "new_user@foo.com",
		Enabled:  false,
		TenantID: "12345",
	}

	th.AssertDeepEquals(t, expected, user)
}
