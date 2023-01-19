//go:build acceptance || compute || images
// +build acceptance compute images

package v2

import (
	"testing"

	"github.com/bizflycloud/gophercloud/acceptance/clients"
	"github.com/bizflycloud/gophercloud/acceptance/tools"
	"github.com/bizflycloud/gophercloud/openstack/compute/v2/images"
	th "github.com/bizflycloud/gophercloud/testhelper"
)

func TestImagesList(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	allPages, err := images.ListDetail(client, nil).AllPages()
	th.AssertNoErr(t, err)

	allImages, err := images.ExtractImages(allPages)
	th.AssertNoErr(t, err)

	var found bool
	for _, image := range allImages {
		tools.PrintResource(t, image)

		if image.ID == choices.ImageID {
			found = true
		}
	}

	th.AssertEquals(t, found, true)
}

func TestImagesGet(t *testing.T) {
	client, err := clients.NewComputeV2Client()
	th.AssertNoErr(t, err)

	choices, err := clients.AcceptanceTestChoicesFromEnv()
	th.AssertNoErr(t, err)

	image, err := images.Get(client, choices.ImageID).Extract()
	th.AssertNoErr(t, err)

	tools.PrintResource(t, image)

	th.AssertEquals(t, choices.ImageID, image.ID)
}
