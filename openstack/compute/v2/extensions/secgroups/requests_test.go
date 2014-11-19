package secgroups

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

const (
	serverID = "{serverID}"
	groupID  = "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5"
)

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

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockGetGroupsResponse(t, groupID)

	group, err := Get(client.ServiceClient(), groupID).Extract()
	th.AssertNoErr(t, err)

	expected := &SecurityGroup{
		ID:          "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
		Description: "default",
		Name:        "default",
		TenantID:    "openstack",
		Rules: []Rule{
			Rule{
				FromPort:      80,
				ToPort:        85,
				IPProtocol:    "TCP",
				IPRange:       IPRange{CIDR: "0.0.0.0"},
				Group:         Group{TenantID: "openstack", Name: "default"},
				ParentGroupID: "b0e0d7dd-2ca4-49a9-ba82-c44a148b66a5",
				ID:            "ebe599e2-6b8c-457c-b1ff-a75e48f10923",
			},
		},
	}

	th.AssertDeepEquals(t, expected, group)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	mockDeleteGroupResponse(t, groupID)

	err := Delete(client.ServiceClient(), groupID).ExtractErr()
	th.AssertNoErr(t, err)
}
