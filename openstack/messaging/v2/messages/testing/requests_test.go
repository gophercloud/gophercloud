package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/messaging/v2/messages"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleListSuccessfully(t, fakeServer)

	listOpts := messages.ListOpts{
		Limit: 1,
	}

	count := 0
	err := messages.List(client.ServiceClient(fakeServer), QueueName, listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
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
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleCreateSuccessfully(t, fakeServer)

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

	actual, err := messages.Create(context.TODO(), client.ServiceClient(fakeServer), QueueName, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedResources, actual)
}

func TestGetMessages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetMessagesSuccessfully(t, fakeServer)

	getMessagesOpts := messages.GetMessagesOpts{
		IDs: []string{"9988776655"},
	}

	actual, err := messages.GetMessages(context.TODO(), client.ServiceClient(fakeServer), QueueName, getMessagesOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedMessagesSet, actual)
}

func TestGetMessage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleGetSuccessfully(t, fakeServer)

	actual, err := messages.Get(context.TODO(), client.ServiceClient(fakeServer), QueueName, MessageID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, FirstMessage, actual)
}

func TestDeleteMessages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteMessagesSuccessfully(t, fakeServer)

	deleteMessagesOpts := messages.DeleteMessagesOpts{
		IDs: []string{"9988776655"},
	}

	err := messages.DeleteMessages(context.TODO(), client.ServiceClient(fakeServer), QueueName, deleteMessagesOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestPopMessages(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandlePopSuccessfully(t, fakeServer)

	popMessagesOpts := messages.PopMessagesOpts{
		Pop: 1,
	}

	actual, err := messages.PopMessages(context.TODO(), client.ServiceClient(fakeServer), QueueName, popMessagesOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedPopMessage, actual)
}

func TestDeleteMessage(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()
	HandleDeleteSuccessfully(t, fakeServer)

	deleteOpts := messages.DeleteOpts{
		ClaimID: "12345",
	}

	err := messages.Delete(context.TODO(), client.ServiceClient(fakeServer), QueueName, MessageID, deleteOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
