// +build acceptance messaging claims

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
)

func TestCRUDClaim(t *testing.T) {
	clientID := "3381af92-2b9e-11e3-b191-71861300734c"

	client, err := clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}

	createdQueueName, err := CreateQueue(t, client)
	defer DeleteQueue(t, client, createdQueueName)

	clientID = "3381af92-2b9e-11e3-b191-71861300734d"

	client, err = clients.NewMessagingV2Client(clientID)
	if err != nil {
		t.Fatalf("Unable to create a messaging service client: %v", err)
	}
	for i := 0; i < 3; i++ {
		CreateMessage(t, client, createdQueueName)
	}

	clientID = "3381af92-2b9e-11e3-b191-7186130073dd"
	claimedMessages, err := CreateClaim(t, client, createdQueueName)

	tools.PrintResource(t, claimedMessages)
}
