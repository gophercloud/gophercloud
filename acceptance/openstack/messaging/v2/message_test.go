// +build acceptance messaging messages

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/messaging/v2/messages"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestListMessages(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-718613007343"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	createdQueueName, err := CreateQueue(t, client)
	defer DeleteQueue(t, client, createdQueueName)

	for i := 0; i < 3; i++ {
		CreateMessage(t, client, createdQueueName)
	}

	// Use a different client/clientID in order to see messages on the Queue
	clientID = "3381af92-2b9e-11e3-b191-71861300734d"
	client, err = clients.NewMessagingV2Client(clientID)

	listOpts := messages.ListOpts{}

	pager := messages.List(client, createdQueueName, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		allMessages, err := messages.ExtractMessages(page)
		if err != nil {
			t.Fatalf("Unable to extract messages: %v", err)
		}

		for _, message := range allMessages {
			tools.PrintResource(t, message)
		}

		return true, nil
	})
}

func TestCreateMessages(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-71861300734c"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	createdQueueName, err := CreateQueue(t, client)
	defer DeleteQueue(t, client, createdQueueName)

	CreateMessage(t, client, createdQueueName)
}
