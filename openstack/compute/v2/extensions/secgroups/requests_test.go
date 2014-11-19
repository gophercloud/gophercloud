package secgroups

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

const serverID = "{serverID}"

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockListGroupsResponse(t)

	count := 0

	err := List(client.ServiceClient()).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractSecurityGroups(page)
		if err != nil {
			t.Errorf("Failed to extract users: %v", err)
			return false, err
		}

		expected := []SecurityGroup{
			SecurityGroup{
				ID:          "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
				Description: "default",
				Name:        "default",
				Rules:       []Rule{},
				TenantID:    "openstack",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestListByServer(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockListGroupsByServerResponse(t, serverID)

	count := 0

	err := ListByServer(client.ServiceClient(), serverID).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractSecurityGroups(page)
		if err != nil {
			t.Errorf("Failed to extract users: %v", err)
			return false, err
		}

		expected := []SecurityGroup{
			SecurityGroup{
				ID:          "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
				Description: "default",
				Name:        "default",
				Rules:       []Rule{},
				TenantID:    "openstack",
			},
		}

		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})

	th.AssertNoErr(t, err)
	th.AssertEquals(t, 1, count)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockCreateGroupResponse(t)

	opts := CreateOpts{
		Name:        "test",
		Description: "something",
	}

	group, err := Create(client.ServiceClient(), opts).Extract()
	th.AssertNoErr(t, err)

	expected := &SecurityGroup{
		ID:          "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
		Name:        "test",
		Description: "something",
		TenantID:    "openstack",
		Rules:       []Rule{},
	}
	th.AssertDeepEquals(t, expected, group)
}
