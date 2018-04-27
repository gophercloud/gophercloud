package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/messaging/v2/messages"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	listOpts := messages.ListOpts{
		Limit: 1,
	}

	count := 0
	err := messages.List(fake.ServiceClient(), QueueName, listOpts).EachPage(func(page pagination.Page) (bool, error) {
		actual, err := messages.ExtractMessages(page)
		th.AssertNoErr(t, err)

		th.CheckDeepEquals(t, ExpectedMessagesSlice[count], actual)
		count++

		return true, nil
	})
	th.AssertNoErr(t, err)

	th.CheckEquals(t, 2, count)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)

	createOpts := messages.BatchCreateOpts{
		messages.CreateOpts{
			TTL:   300,
			Delay: 20,
			Body: map[string]interface{}{
				"event":     "BackupStarted",
				"backup_id": "c378813c-3f0b-11e2-ad92-7823d2b0f3ce",
			},
		},
		messages.CreateOpts{
			Body: map[string]interface{}{
				"event":         "BackupProgress",
				"current_bytes": "0",
				"total_bytes":   "99614720",
			},
		},
	}

	actual, err := messages.Create(fake.ServiceClient(), QueueName, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedResources, actual)
}
