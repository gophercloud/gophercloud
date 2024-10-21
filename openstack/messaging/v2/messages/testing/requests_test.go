package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/messaging/v2/messages"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	listOpts := messages.ListOpts{
		Limit: 1,
	}

	count := 0
	err := messages.List(fake.ServiceClient(), QueueName, listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
			Body: map[string]any{
				"event":     "BackupStarted",
				"backup_id": "c378813c-3f0b-11e2-ad92-7823d2b0f3ce",
			},
		},
		messages.CreateOpts{
			Body: map[string]any{
				"event":         "BackupProgress",
				"current_bytes": "0",
				"total_bytes":   "99614720",
			},
		},
	}

	actual, err := messages.Create(context.TODO(), fake.ServiceClient(), QueueName, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedResources, actual)
}

func TestGetMessages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetMessagesSuccessfully(t)

	getMessagesOpts := messages.GetMessagesOpts{
		IDs: []string{"9988776655"},
	}

	actual, err := messages.GetMessages(context.TODO(), fake.ServiceClient(), QueueName, getMessagesOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedMessagesSet, actual)
}

func TestGetMessage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	actual, err := messages.Get(context.TODO(), fake.ServiceClient(), QueueName, MessageID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, FirstMessage, actual)
}

func TestDeleteMessages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteMessagesSuccessfully(t)

	deleteMessagesOpts := messages.DeleteMessagesOpts{
		IDs: []string{"9988776655"},
	}

	err := messages.DeleteMessages(context.TODO(), fake.ServiceClient(), QueueName, deleteMessagesOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestPopMessages(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePopSuccessfully(t)

	popMessagesOpts := messages.PopMessagesOpts{
		Pop: 1,
	}

	actual, err := messages.PopMessages(context.TODO(), fake.ServiceClient(), QueueName, popMessagesOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedPopMessage, actual)
}

func TestDeleteMessage(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)

	deleteOpts := messages.DeleteOpts{
		ClaimID: "12345",
	}

	err := messages.Delete(context.TODO(), fake.ServiceClient(), QueueName, MessageID, deleteOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
