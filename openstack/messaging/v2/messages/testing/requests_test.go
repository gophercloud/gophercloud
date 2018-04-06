package testing

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/messaging/v2/messages"
	th "github.com/gophercloud/gophercloud/testhelper"
	fake "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)

	createOpts := messages.BatchCreateOpts{
		messages.CreateOpts{
			TTL:   300,
			Delay: 20,
			Body: map[string]interface{}{
				"event":     "BackupStarted",
				"backup_id": "c378813c-3f0b-11e2-ad92-7823d2b0f3ce",
			},
		},
		messages.CreateOpts{
			Body: map[string]interface{}{
				"event":         "BackupProgress",
				"current_bytes": "0",
				"total_bytes":   "99614720",
			},
		},
	}

	actual, err := messages.Create(fake.ServiceClient(), QueueName, createOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedResources, actual)
}
