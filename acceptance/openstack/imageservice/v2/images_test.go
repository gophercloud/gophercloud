// +build acceptance imageservice images

package v2

import (
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/acceptance/clients"
	"github.com/gophercloud/gophercloud/acceptance/tools"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestImagesListEachPage(t *testing.T) {
	client, err := clients.NewImageServiceV2Client()
	if err != nil {
		t.Fatalf("Unable to create an image service client: %v", err)
	}

	listOpts := images.ListOpts{
		Limit: 1,
	}

	pager := images.List(client, listOpts)
	err = pager.EachPage(func(page pagination.Page) (bool, error) {
		images, err := images.ExtractImages(page)
		if err != nil {
			t.Fatalf("Unable to extract images: %v", err)
		}

		for _, image := range images {
			tools.PrintResource(t, image)
			tools.PrintResource(t, image.Properties)
		}

		return true, nil
	})
}

func TestImagesListAllPages(t *testing.T) {
	client, err := clients.NewImageServiceV2Client()
	if err != nil {
		t.Fatalf("Unable to create an image service client: %v", err)
	}

	listOpts := images.ListOpts{
		Limit: 1,
	}

	allPages, err := images.List(client, listOpts).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve all images: %v", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		t.Fatalf("Unable to extract images: %v", err)
	}

	for _, image := range allImages {
		tools.PrintResource(t, image)
		tools.PrintResource(t, image.Properties)
	}
}

func TestImagesCreateDestroyEmptyImage(t *testing.T) {
	client, err := clients.NewImageServiceV2Client()
	if err != nil {
		t.Fatalf("Unable to create an image service client: %v", err)
	}

	image, err := CreateEmptyImage(t, client)
	if err != nil {
		t.Fatalf("Unable to create empty image: %v", err)
	}

	defer DeleteImage(t, client, image)

	tools.PrintResource(t, image)
}

func TestImagesListByDate(t *testing.T) {
	client, err := clients.NewImageServiceV2Client()
	if err != nil {
		t.Fatalf("Unable to create an image service client: %v", err)
	}

	date := time.Date(2014, 1, 1, 1, 1, 1, 0, time.UTC)
	listOpts := images.ListOpts{
		Limit: 1,
		CreatedAt: &images.ImageDateQuery{
			Date:   date,
			Filter: images.FilterGTE,
		},
	}

	allPages, err := images.List(client, listOpts).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve all images: %v", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		t.Fatalf("Unable to extract images: %v", err)
	}

	for _, image := range allImages {
		tools.PrintResource(t, image)
		tools.PrintResource(t, image.Properties)
	}

	date = time.Date(2049, 1, 1, 1, 1, 1, 0, time.UTC)
	listOpts = images.ListOpts{
		Limit: 1,
		CreatedAt: &images.ImageDateQuery{
			Date:   date,
			Filter: images.FilterGTE,
		},
	}

	allPages, err = images.List(client, listOpts).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve all images: %v", err)
	}

	allImages, err = images.ExtractImages(allPages)
	if err != nil {
		t.Fatalf("Unable to extract images: %v", err)
	}

	if len(allImages) > 0 {
		t.Fatalf("Expected 0 images, got %d", len(allImages))
	}
}
