// +build acceptance messaging messages

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/acceptance/clients"
)

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
