package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/messaging/v2/queues"
)

func CreateQueue(t *testing.T, client *gophercloud.ServiceClient) (string, error) {
	queueName := tools.RandomString("ACPTTEST", 5)

	t.Logf("Attempting to create Queue: %s", queueName)

	createOpts := queues.CreateOpts{
		QueueName:                  queueName,
		MaxMessagesPostSize:        262143,
		DefaultMessageTTL:          3700,
		DefaultMessageDelay:        25,
		DeadLetterQueueMessagesTTL: 3500,
		MaxClaimCount:              10,
		Extra:                      map[string]interface{}{"description": "Test Queue for Gophercloud acceptance tests."},
	}

	createErr := queues.Create(client, createOpts).ExtractErr()
	if createErr != nil {
		t.Fatalf("Unable to create Queue: %v", createErr)
	}

	GetQueue(t, client, queueName)

	t.Logf("Created Queue: %s", queueName)
	return queueName, nil
}

func DeleteQueue(t *testing.T, client *gophercloud.ServiceClient, queueName string) {
	t.Logf("Attempting to delete Queue: %s", queueName)
	err := queues.Delete(client, queueName).ExtractErr()
	if err != nil {
		t.Fatalf("Unable to delete Queue %s: %v", queueName, err)
	}

	t.Logf("Deleted Queue: %s", queueName)
}

func GetQueue(t *testing.T, client *gophercloud.ServiceClient, queueName string) (queues.QueueDetails, error) {
	t.Logf("Attempting to get Queue: %s", queueName)
	queue, err := queues.Get(client, queueName).Extract()
	if err != nil {
		t.Fatalf("Unable to get Queue %s: %v", queueName, err)
	}
	return queue, nil
}
