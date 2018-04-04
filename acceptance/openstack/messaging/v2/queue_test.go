package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/messaging/v2/queues"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
)

func TestCRUDQueues(t *testing.T) {
	client, err := clients.NewMessagingV2Client()
	th.AssertNoErr(t, err)

	clientID := "3381af92-2b9e-11e3-b191-71861300734d"

	createdQueueName, err := CreateQueue(t, client, clientID)
	defer DeleteQueue(t, client, createdQueueName, clientID)

	createdQueue, err := queues.Get(client, createdQueueName, clientID).Extract()
	tools.PrintResource(t, createdQueue)

	updateOpts := queues.UpdateOpts{
		queues.UpdateQueueBody{
			Op:    "replace",
			Path:  "/metadata/_max_claim_count",
			Value: 15,
		},
		queues.UpdateQueueBody{
			Op:    "replace",
			Path:  "/metadata/description",
			Value: "Updated description for queues acceptance test.",
		},
	}

	updateResult, updateErr := queues.Update(client, createdQueueName, clientID, updateOpts).Extract()
	th.AssertNoErr(t, updateErr)

	tools.PrintResource(t, updateResult)
}

func TestListQueues(t *testing.T) {
	client, err := clients.NewMessagingV2Client()
	th.AssertNoErr(t, err)

	clientID := "3381af92-2b9e-11e3-b191-71861300734d"

	firstQueueName, err := CreateQueue(t, client, clientID)
	defer DeleteQueue(t, client, firstQueueName, clientID)

	secondQueueName, err := CreateQueue(t, client, clientID)
	defer DeleteQueue(t, client, secondQueueName, clientID)

	listOpts := queues.ListOpts{
		Limit:    10,
		Detailed: true,
	}

	pager := queues.List(client, clientID, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		allQueues, err := queues.ExtractQueues(page)
		th.AssertNoErr(t, err)

		for _, queue := range allQueues {
			tools.PrintResource(t, queue)
		}

		return true, nil
	})
}
