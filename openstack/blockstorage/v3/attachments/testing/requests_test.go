package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/attachments"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

func TestListAll(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockListResponse(t)

	allPages, err := attachments.List(client.ServiceClient(), &attachments.ListOpts{}).AllPages(context.TODO())
	th.AssertNoErr(t, err)
	actual, err := attachments.ExtractAttachments(allPages)
	th.AssertNoErr(t, err)

	expected := []attachments.Attachment{*expectedAttachment}

	th.CheckDeepEquals(t, expected, actual)

}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockGetResponse(t)

	attachment, err := attachments.Get(context.TODO(), client.ServiceClient(), "05551600-a936-4d4a-ba42-79a037c1-c91a").Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, expectedAttachment, attachment)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCreateResponse(t)

	options := &attachments.CreateOpts{
		InstanceUUID: "83ec2e3b-4321-422b-8706-a84185f52a0a",
		Connector: map[string]any{
			"initiator":  "iqn.1993-08.org.debian: 01: cad181614cec",
			"ip":         "192.168.1.20",
			"platform":   "x86_64",
			"host":       "tempest-1",
			"os_type":    "linux2",
			"multipath":  false,
			"mountpoint": "/dev/vdb",
			"mode":       "rw",
		},
		VolumeUUID: "289da7f8-6440-407c-9fb4-7db01ec49164",
	}
	attachment, err := attachments.Create(context.TODO(), client.ServiceClient(), options).Extract()
	th.AssertNoErr(t, err)

	th.AssertDeepEquals(t, expectedAttachment, attachment)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockDeleteResponse(t)

	res := attachments.Delete(context.TODO(), client.ServiceClient(), "05551600-a936-4d4a-ba42-79a037c1-c91a")
	th.AssertNoErr(t, res.Err)
}

func TestUpdate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateResponse(t)

	options := &attachments.UpdateOpts{
		Connector: map[string]any{
			"initiator":  "iqn.1993-08.org.debian: 01: cad181614cec",
			"ip":         "192.168.1.20",
			"platform":   "x86_64",
			"host":       "tempest-1",
			"os_type":    "linux2",
			"multipath":  false,
			"mountpoint": "/dev/vdb",
			"mode":       "rw",
		},
	}
	attachment, err := attachments.Update(context.TODO(), client.ServiceClient(), "05551600-a936-4d4a-ba42-79a037c1-c91a", options).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expectedAttachment, attachment)
}

func TestUpdateEmpty(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockUpdateEmptyResponse(t)

	options := attachments.UpdateOpts{}
	attachment, err := attachments.Update(context.TODO(), client.ServiceClient(), "05551600-a936-4d4a-ba42-79a037c1-c91a", options).Extract()
	th.AssertNoErr(t, err)
	th.AssertDeepEquals(t, expectedAttachment, attachment)
}

func TestComplete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockCompleteResponse(t)

	err := attachments.Complete(context.TODO(), client.ServiceClient(), "05551600-a936-4d4a-ba42-79a037c1-c91a").ExtractErr()
	th.AssertNoErr(t, err)
}
