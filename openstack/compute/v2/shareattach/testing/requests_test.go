package testing

import (
	"context"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/shareattach"
	"github.com/gophercloud/gophercloud/v2/pagination"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	"github.com/gophercloud/gophercloud/v2/testhelper/client"
)

var adminExportLocation = "10.0.0.50:/mnt/foo"
var adminUUID = "68ba1762-fd6d-4221-8311-f3193dd93404"

// FirstShareAttachment is the first result in ListOutput.
var FirstShareAttachment = shareattach.ShareAttachment{
	ShareID: "e8debdc0-447a-4376-a10a-4cd9122d7986",
	Status:  "inactive",
	Tag:     "e8debdc0-447a-4376-a10a-4cd9122d7986",
}

// SecondShareAttachment is the 2nd result in ListOutput.
var SecondShareAttachment = shareattach.ShareAttachment{
	ShareID: "a26887c6-c47b-4654-abb5-dfadf7d3f803",
	Status:  "active",
	Tag:     "a26887c6-c47b-4654-abb5-dfadf7d3f803",
}

// FirstShareAttachmentAdmin is the parsed result from GetOutputAdmin.
var FirstShareAttachmentAdmin = shareattach.ShareAttachment{
	ShareID:        "e8debdc0-447a-4376-a10a-4cd9122d7986",
	Status:         "inactive",
	Tag:            "e8debdc0-447a-4376-a10a-4cd9122d7986",
	ExportLocation: &adminExportLocation,
	UUID:           &adminUUID,
}

// ExpectedShareAttachmentSlice is the slice of results that should be parsed
// from ListOutput, in the expected order.
var ExpectedShareAttachmentSlice = []shareattach.ShareAttachment{FirstShareAttachment, SecondShareAttachment}

// CreatedShareAttachment is the parsed result from CreateOutput.
var CreatedShareAttachment = shareattach.ShareAttachment{
	ShareID: "e8debdc0-447a-4376-a10a-4cd9122d7986",
	Status:  "attaching",
	Tag:     "my-share",
}

// CreatedShareAttachmentWithoutTag is the parsed result from CreateOutputWithoutTag.
// Nova returns tag equal to share_id when the client omits tag on create.
var CreatedShareAttachmentWithoutTag = shareattach.ShareAttachment{
	ShareID: "e8debdc0-447a-4376-a10a-4cd9122d7986",
	Status:  "attaching",
	Tag:     "e8debdc0-447a-4376-a10a-4cd9122d7986",
}

func TestList(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleListSuccessfully(t, fakeServer)

	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	count := 0
	err := shareattach.List(client.ServiceClient(fakeServer), serverID).EachPage(context.TODO(), func(_ context.Context, page pagination.Page) (bool, error) {
		count++
		actual, err := shareattach.ExtractShareAttachments(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedShareAttachmentSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestCreate(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreateSuccessfully(t, fakeServer)

	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	actual, err := shareattach.Create(context.TODO(), client.ServiceClient(fakeServer), serverID, shareattach.CreateOpts{
		ShareID: "3cdf5132-64f2-4584-876a-bd296ae7eabd",
		Tag:     "my-share",
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedShareAttachment, actual)
}

func TestCreateWithoutTag(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleCreateSuccessfullyWithoutTag(t, fakeServer)

	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	actual, err := shareattach.Create(context.TODO(), client.ServiceClient(fakeServer), serverID, shareattach.CreateOpts{
		ShareID: "3cdf5132-64f2-4584-876a-bd296ae7eabd",
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedShareAttachmentWithoutTag, actual)
}

func TestGet(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetSuccessfully(t, fakeServer)

	shareID := "e8debdc0-447a-4376-a10a-4cd9122d7986"
	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	actual, err := shareattach.Get(context.TODO(), client.ServiceClient(fakeServer), serverID, shareID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstShareAttachment, actual)
}

func TestGetAdmin(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleGetSuccessfullyAdmin(t, fakeServer)

	shareID := "e8debdc0-447a-4376-a10a-4cd9122d7986"
	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	actual, err := shareattach.Get(context.TODO(), client.ServiceClient(fakeServer), serverID, shareID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &FirstShareAttachmentAdmin, actual)
}

func TestDelete(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	HandleDeleteSuccessfully(t, fakeServer)

	shareID := "e8debdc0-447a-4376-a10a-4cd9122d7986"
	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	err := shareattach.Delete(context.TODO(), client.ServiceClient(fakeServer), serverID, shareID).ExtractErr()
	th.AssertNoErr(t, err)
}
