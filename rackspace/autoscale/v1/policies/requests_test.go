package policies

import (
	"testing"

	"github.com/rackspace/gophercloud/pagination"
	th "github.com/rackspace/gophercloud/testhelper"
	"github.com/rackspace/gophercloud/testhelper/client"
)

const (
	groupID         = "60b15dad-5ea1-43fa-9a12-a1d737b4da07"
	webhookPolicyID = "2b48d247-0282-4b9d-8775-5c4b67e8e649"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePolicyListSuccessfully(t)

	pages := 0
	pager := List(client.ServiceClient(), "60b15dad-5ea1-43fa-9a12-a1d737b4da07")

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

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePolicyGetSuccessfully(t)

	client := client.ServiceClient()

	policy, err := Get(client, groupID, webhookPolicyID).Extract()

	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, WebhookPolicy, *policy)
}
