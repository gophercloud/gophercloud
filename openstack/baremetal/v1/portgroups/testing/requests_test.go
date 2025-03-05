package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/v1/portgroups"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"

	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestListPortGroups(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePortGroupListSuccessfully(t)

	pages := 0
	err := portgroups.List(client.ServiceClient(), portgroups.ListOpts{}).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		pages++

		actual, err := portgroups.ExtractPortGroups(page)
		if err != nil {
			return false, err
		}

		if len(actual) != 2 {
			t.Fatalf("Expected 2 portgroups, got %d", len(actual))
		}
		th.AssertEquals(t, "d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796", actual[0].UUID)
		th.AssertEquals(t, "a1b2c3d4-e5f6-7890-1234-56789abcdef0", actual[1].UUID)
		th.AssertEquals(t, "bond0", actual[0].Name)
		th.AssertEquals(t, "bond1", actual[1].Name)

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}

func TestCreatePortGroup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePortGroupCreationSuccessfully(t, SinglePortGroupBody)

	actual, err := portgroups.Create(context.TODO(), client.ServiceClient(), portgroups.CreateOpts{
		Name:                     "bond0",
		NodeUUID:                 "f9c9a846-c53f-4b17-9f0c-dd9f459d35c8",
		Address:                  "00:1a:2b:3c:4d:5e",
		Mode:                     "active-backup",
		StandalonePortsSupported: true,
		Properties: map[string]interface{}{
			"miimon":           "100",
			"updelay":          "1000",
			"downdelay":        "1000",
			"xmit_hash_policy": "layer2",
		},
		Extra: map[string]string{
			"description": "Primary network bond",
			"location":    "rack-3-unit-12",
		},
	}).Extract()
	th.AssertNoErr(t, err)

	th.CheckDeepEquals(t, PortGroup1, *actual)
}

func TestDeletePortGroup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePortGroupDeletionSuccessfully(t)

	res := portgroups.Delete(context.TODO(), client.ServiceClient(), "d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796")
	th.AssertNoErr(t, res.Err)
}

func TestGetPortGroup(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePortGroupGetSuccessfully(t)

	c := client.ServiceClient()
	actual, err := portgroups.Get(context.TODO(), c, "d2b42f0d-c7e6-4f08-b9bc-e8b23a6ee796").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, PortGroup1, *actual)
}
