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
