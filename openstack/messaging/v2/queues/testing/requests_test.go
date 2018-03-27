package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/messaging/v2/queues"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)

	createOpts := queues.CreateOpts{
		MaxMessagesPostSize:       262144,
		DefaultMessageTTL:         3600,
		DefaultMessageDelay:       30,
		DeadLetterQueue:           "dead_letter",
		DeadLetterQueueMessageTTL: 3600,
		MaxClaimCount:             10,
		Description:               "Queue for unit testing.",
	}

	err := queues.Create(fake.ServiceClient(), QueueName, ClientID, createOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
