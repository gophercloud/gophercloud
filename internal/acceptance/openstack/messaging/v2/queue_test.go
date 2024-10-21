//go:build acceptance || messaging || queues

package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/internal/acceptance/clients"
	"github.com/gophercloud/gophercloud/v2/internal/acceptance/tools"
	"github.com/gophercloud/gophercloud/v2/openstack/messaging/v2/queues"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
)

func TestCRUDQueues(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-71861300734d"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	createdQueueName, err := CreateQueue(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQueue(t, client, createdQueueName)

	createdQueue, err := queues.Get(context.TODO(), client, createdQueueName).Extract()
	th.AssertNoErr(t, err)

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
	updateResult, updateErr := queues.Update(context.TODO(), client, createdQueueName, updateOpts).Extract()
	if updateErr != nil {
		t.Fatalf("Unable to update Queue %s: %v", createdQueueName, updateErr)
	}

	updatedQueue, err := GetQueue(t, client, createdQueueName)
	th.AssertNoErr(t, err)

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
	th.AssertNoErr(t, err)
	defer DeleteQueue(t, client, firstQueueName)

	secondQueueName, err := CreateQueue(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQueue(t, client, secondQueueName)

	listOpts := queues.ListOpts{
		Limit:    10,
		Detailed: true,
	}

	pager := queues.List(client, listOpts)
	err = pager.EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		allQueues, err := queues.ExtractQueues(page)
		if err != nil {
			t.Fatalf("Unable to extract Queues: %v", err)
		}

		for _, queue := range allQueues {
			tools.PrintResource(t, queue)
		}

		return true, nil
	})
	th.AssertNoErr(t, err)
}

func TestStatQueue(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-71861300734c"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	createdQueueName, err := CreateQueue(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQueue(t, client, createdQueueName)

	queueStats, err := queues.GetStats(context.TODO(), client, createdQueueName).Extract()
	if err != nil {
		t.Fatalf("Unable to stat queue: %v", err)
	}

	tools.PrintResource(t, queueStats)
}

func TestShare(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-71861300734c"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	queueName, err := CreateQueue(t, client)
	if err != nil {
		t.Logf("Unable to create queue for share.")
	}
	defer DeleteQueue(t, client, queueName)

	t.Logf("Attempting to create share for queue: %s", queueName)
	share, shareErr := CreateShare(t, client, queueName)
	if shareErr != nil {
		t.Fatalf("Unable to create share: %v", shareErr)
	}

	tools.PrintResource(t, share)
}

func TestPurge(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-71861300734c"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	queueName, err := CreateQueue(t, client)
	th.AssertNoErr(t, err)
	defer DeleteQueue(t, client, queueName)

	purgeOpts := queues.PurgeOpts{
		ResourceTypes: []queues.PurgeResource{
			queues.ResourceMessages,
		},
	}

	t.Logf("Attempting to purge queue: %s", queueName)
	purgeErr := queues.Purge(context.TODO(), client, queueName, purgeOpts).ExtractErr()
	if purgeErr != nil {
		t.Fatalf("Unable to purge queue %s: %v", queueName, purgeErr)
	}
}
