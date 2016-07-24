// +build acceptance compute images

package v2

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/images"
)

func TestImagesList(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute: client: %v", err)
	}

	allPages, err := images.ListDetail(client, nil).AllPages()
	if err != nil {
		t.Fatalf("Unable to retrieve images: %v", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		t.Fatalf("Unable to extract image results: %v", err)
	}

	for _, image := range allImages {
		printImage(t, image)
	}
}

func TestImagesGet(t *testing.T) {
	client, err := newClient()
	if err != nil {
		t.Fatalf("Unable to create a compute: client: %v", err)
	}

	choices, err := ComputeChoicesFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	image, err := images.Get(client, choices.ImageID).Extract()
	if err != nil {
		t.Fatalf("Unable to get image information: %v", err)
	}

	printImage(t, *image)
}

func printImage(t *testing.T, image images.Image) {
	t.Logf("ID: %s", image.ID)
	t.Logf("Name: %s", image.Name)
	t.Logf("MinDisk: %d", image.MinDisk)
	t.Logf("MinRAM: %d", image.MinRAM)
	t.Logf("Status: %s", image.Status)
	t.Logf("Progress: %d", image.Progress)
	t.Logf("Metadata: %#v", image.Metadata)
	t.Logf("Created: %s", image.Created)
	t.Logf("Updated: %s", image.Updated)
}
