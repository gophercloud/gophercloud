package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/messaging/v2/queues"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestCRUDQueues(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-71861300734d"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	createdQueueName, err := CreateQueue(t, client)
	defer DeleteQueue(t, client, createdQueueName)

	createdQueue, err := queues.Get(client, createdQueueName).Extract()

	tools.PrintResource(t, createdQueue)
	tools.PrintResource(t, createdQueue.Extra)

	updateOpts := queues.BatchUpdateOpts{
		queues.UpdateOpts{
			Op:    "replace",
			Path:  "/metadata/_max_claim_count",
			Value: 15,
		},
		queues.UpdateOpts{
			Op:    "replace",
			Path:  "/metadata/description",
			Value: "Updated description for queues acceptance test.",
		},
	}

	t.Logf("Attempting to update Queue: %s", createdQueueName)
	updateResult, updateErr := queues.Update(client, createdQueueName, updateOpts).Extract()
	if updateErr != nil {
		t.Fatalf("Unable to update Queue %s: %v", createdQueueName, updateErr)
	}

	updatedQueue, err := GetQueue(t, client, createdQueueName)

	tools.PrintResource(t, updateResult)
	tools.PrintResource(t, updatedQueue)
	tools.PrintResource(t, updatedQueue.Extra)
}

func TestListQueues(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-71861300734d"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	firstQueueName, err := CreateQueue(t, client)
	defer DeleteQueue(t, client, firstQueueName)

	secondQueueName, err := CreateQueue(t, client)
	defer DeleteQueue(t, client, secondQueueName)

	listOpts := queues.ListOpts{
		Limit:    10,
		Detailed: true,
	}

	pager := queues.List(client, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		allQueues, err := queues.ExtractQueues(page)
		if err != nil {
			t.Fatalf("Unable to extract Queues: %v", err)
		}

		for _, queue := range allQueues {
			tools.PrintResource(t, queue)
		}

		return true, nil
	})
}
