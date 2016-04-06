package policies

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

const (
	groupID = "10eb3219-1b12-4b34-b1e4-e10ee4f24c65"
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

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePolicyCreateSuccessfully(t)

	client := client.ServiceClient()
	opts := CreateOpts{
		{
			Name:     "webhook policy",
			Type:     Webhook,
			Cooldown: 300,
			Adjustment: Adjustment{
				Type:  ChangePercent,
				Value: 3.3,
			},
		},
		{
			Name: "one time",
			Type: Schedule,
			Adjustment: Adjustment{
				Type:  Change,
				Value: -1,
			},
			Args: map[string]interface{}{
				"at": "2020-04-01T23:00:00.000Z",
			},
		},
		{
			Name: "sunday afternoon",
			Type: Schedule,
			Adjustment: Adjustment{
				Type:  DesiredCapacity,
				Value: 2,
			},
			Args: map[string]interface{}{
				"cron": "59 15 * * 0",
			},
		},
	}

	policies, err := Create(client, groupID, opts).Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, WebhookPolicy, policies[0])
	th.CheckDeepEquals(t, OneTimePolicy, policies[1])
	th.CheckDeepEquals(t, SundayAfternoonPolicy, policies[2])
}
