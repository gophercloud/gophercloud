package policies

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePolicyListSuccessfully(t)

	pages := 0
	pager := List(client.ServiceClient(), "10eb3219-1b12-4b34-b1e4-e10ee4f24c65")

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		pages++

		policies, err := ExtractPolicies(page)

		if err != nil {
			return false, err
		}

		if len(policies) != 3 {
			t.Fatalf("Expected 3 policies, got %d", len(policies))
		}

		th.CheckDeepEquals(t, WebhookPolicy, policies[0])
		th.CheckDeepEquals(t, OneTimePolicy, policies[1])
		th.CheckDeepEquals(t, SundayAfternoonPolicy, policies[2])

		return true, nil
	})

	th.AssertNoErr(t, err)

	if pages != 1 {
		t.Errorf("Expected 1 page, saw %d", pages)
	}
}
