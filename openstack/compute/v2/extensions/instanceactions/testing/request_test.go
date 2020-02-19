package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/instanceactions"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleInstanceActionListSuccessfully(t)

	expected := ListExpected
	pages := 0
	err := instanceactions.List(client.ServiceClient(), "asdfasdfasdf", nil).EachPage(func(page pagination.Page) (bool, error) {
		pages++

		actual, err := instanceactions.ExtractInstanceActions(page)
		th.AssertNoErr(t, err)

		if len(actual) != 2 {
			t.Fatalf("Expected 2 instance actions, got %d", len(actual))
		}
		th.CheckDeepEquals(t, expected, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, pages)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleInstanceActionGetSuccessfully(t)

	client := client.ServiceClient()
	actual, err := instanceactions.Get(client, "asdfasdfasdf", "okzeorkmkfs").Extract()
	if err != nil {
		t.Fatalf("Unexpected Get error: %v", err)
	}

	th.CheckDeepEquals(t, GetExpected, actual)
}
