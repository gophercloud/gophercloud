package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/messaging/v2/queues"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	listOpts := queues.ListOpts{}

	count := 0
	err := queues.List(fake.ServiceClient(), ClientID, listOpts).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := queues.ExtractQueues(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedQueueSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)

	createOpts := queues.CreateOpts{
		QueueName:                  QueueName,
		MaxMessagesPostSize:        262144,
		DefaultMessageTTL:          3600,
		DefaultMessageDelay:        30,
		DeadLetterQueue:            "dead_letter",
		DeadLetterQueueMessagesTTL: 3600,
		MaxClaimCount:              10,
		Extra:                      map[string]interface{}{"description": "Queue for unit testing."},
	}

	err := queues.Create(fake.ServiceClient(), ClientID, createOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
