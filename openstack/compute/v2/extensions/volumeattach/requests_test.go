package volumeattach

import (
	"testing"

	fixtures "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach/testing"
	"github.com/gophercloud/gophercloud/pagination"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// FirstVolumeAttachment is the first result in ListOutput.
var FirstVolumeAttachment = VolumeAttachment{
	Device:   "/dev/vdd",
	ID:       "a26887c6-c47b-4654-abb5-dfadf7d3f803",
	ServerID: "4d8c3732-a248-40ed-bebc-539a6ffd25c0",
	VolumeID: "a26887c6-c47b-4654-abb5-dfadf7d3f803",
}

// SecondVolumeAttachment is the first result in ListOutput.
var SecondVolumeAttachment = VolumeAttachment{
	Device:   "/dev/vdc",
	ID:       "a26887c6-c47b-4654-abb5-dfadf7d3f804",
	ServerID: "4d8c3732-a248-40ed-bebc-539a6ffd25c0",
	VolumeID: "a26887c6-c47b-4654-abb5-dfadf7d3f804",
}

// ExpectedVolumeAttachmentSlide is the slice of results that should be parsed
// from ListOutput, in the expected order.
var ExpectedVolumeAttachmentSlice = []VolumeAttachment{FirstVolumeAttachment, SecondVolumeAttachment}

//CreatedVolumeAttachment is the parsed result from CreatedOutput.
var CreatedVolumeAttachment = VolumeAttachment{
	Device:   "/dev/vdc",
	ID:       "a26887c6-c47b-4654-abb5-dfadf7d3f804",
	ServerID: "4d8c3732-a248-40ed-bebc-539a6ffd25c0",
	VolumeID: "a26887c6-c47b-4654-abb5-dfadf7d3f804",
}

func TestList(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixtures.HandleListSuccessfully(t)
	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	count := 0
	err := List(client.ServiceClient(), serverID).EachPage(func(page pagination.Page) (bool, error) {
		count++
		actual, err := ExtractVolumeAttachments(page)
		th.AssertNoErr(t, err)
		th.CheckDeepEquals(t, ExpectedVolumeAttachmentSlice, actual)

		return true, nil
	})
	th.AssertNoErr(t, err)
	th.CheckEquals(t, 1, count)
}

func TestCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixtures.HandleCreateSuccessfully(t)
	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	actual, err := Create(client.ServiceClient(), serverID, CreateOpts{
		Device:   "/dev/vdc",
		VolumeID: "a26887c6-c47b-4654-abb5-dfadf7d3f804",
	}).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &CreatedVolumeAttachment, actual)
}

func TestGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixtures.HandleGetSuccessfully(t)
	aID := "a26887c6-c47b-4654-abb5-dfadf7d3f804"
	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	actual, err := Get(client.ServiceClient(), serverID, aID).Extract()
	th.AssertNoErr(t, err)
	th.CheckDeepEquals(t, &SecondVolumeAttachment, actual)
}

func TestDelete(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	fixtures.HandleDeleteSuccessfully(t)
	aID := "a26887c6-c47b-4654-abb5-dfadf7d3f804"
	serverID := "4d8c3732-a248-40ed-bebc-539a6ffd25c0"

	err := Delete(client.ServiceClient(), serverID, aID).ExtractErr()
	th.AssertNoErr(t, err)
}
