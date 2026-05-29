package v2

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/messages"
)

// DeleteMessage will delete a message. An error will occur if
// the message was unable to be deleted.
func DeleteMessage(t *testing.T, client *gophercloud.ServiceClient, message *messages.Message) {
	err := messages.Delete(context.TODO(), client, message.ID).ExtractErr()
	if err != nil {
		t.Fatalf("Failed to delete message %s: %v", message.ID, err)
	}

	t.Logf("Deleted message: %s", message.ID)
}
