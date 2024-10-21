package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/messaging/v2/queues"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	fake "github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListSuccessfully(t)

	listOpts := queues.ListOpts{
		Limit:     1,
		WithCount: true,
	}

	count := 0
	err := queues.List(fake.ServiceClient(), listOpts).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		actual, err := queues.ExtractQueues(page)
		th.AssertNoErr(t, err)
		countField, err := page.(queues.QueuePage).GetCount()

		th.AssertNoErr(t, err)
		th.AssertEquals(t, countField, 2)

		th.CheckDeepEquals(t, ExpectedQueueSlice[count], actual)
		count++
		return true, nil
	})
	th.AssertNoErr(t, err)

	th.CheckEquals(t, 2, count)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleCreateSuccessfully(t)
	var enableEncrypted *bool = new(bool)

	createOpts := queues.CreateOpts{
		QueueName:                  QueueName,
		MaxMessagesPostSize:        262144,
		DefaultMessageTTL:          3600,
		DefaultMessageDelay:        30,
		DeadLetterQueue:            "dead_letter",
		DeadLetterQueueMessagesTTL: 3600,
		MaxClaimCount:              10,
		EnableEncryptMessages:      enableEncrypted,
		Extra:                      map[string]any{"description": "Queue for unit testing."},
	}

	err := queues.Create(context.TODO(), fake.ServiceClient(), createOpts).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleUpdateSuccessfully(t)

	updateOpts := queues.BatchUpdateOpts{
		queues.UpdateOpts{
			Op:    queues.ReplaceOp,
			Path:  "/metadata/description",
			Value: "Update queue description",
		},
	}
	updatedQueueResult := queues.QueueDetails{
		Extra: map[string]any{"description": "Update queue description"},
	}

	actual, err := queues.Update(context.TODO(), fake.ServiceClient(), QueueName, updateOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, updatedQueueResult, actual)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetSuccessfully(t)

	actual, err := queues.Get(context.TODO(), fake.ServiceClient(), QueueName).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, QueueDetails, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleDeleteSuccessfully(t)

	err := queues.Delete(context.TODO(), fake.ServiceClient(), QueueName).ExtractErr()
	th.AssertNoErr(t, err)
}

func TestGetStat(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleGetStatsSuccessfully(t)

	actual, err := queues.GetStats(context.TODO(), fake.ServiceClient(), QueueName).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedStats, actual)
}

func TestShare(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleShareSuccessfully(t)

	shareOpts := queues.ShareOpts{
		Paths:   []queues.SharePath{queues.PathMessages, queues.PathClaims, queues.PathSubscriptions},
		Methods: []queues.ShareMethod{queues.MethodGet, queues.MethodPost, queues.MethodPut, queues.MethodPatch},
		Expires: "2016-09-01T00:00:00",
	}

	actual, err := queues.Share(context.TODO(), fake.ServiceClient(), QueueName, shareOpts).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, ExpectedShare, actual)
}

func TestPurge(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandlePurgeSuccessfully(t)

	purgeOpts := queues.PurgeOpts{
		ResourceTypes: []queues.PurgeResource{queues.ResourceMessages, queues.ResourceSubscriptions},
	}

	err := queues.Purge(context.TODO(), fake.ServiceClient(), QueueName, purgeOpts).ExtractErr()
	th.AssertNoErr(t, err)
}
