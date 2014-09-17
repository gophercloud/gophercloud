package images

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

// Image is used for JSON (un)marshalling.
// It provides a description of an OS image.
type Image struct {
	// ID contains the image's unique identifier.
	ID string

	Created string

	// MinDisk and MinRAM specify the minimum resources a server must provide to be able to install the image.
	MinDisk int
	MinRAM  int

	// Name provides a human-readable moniker for the OS image.
	Name string

	// The Progress and Status fields indicate image-creation status.
	// Any usable image will have 100% progress.
	Progress int
	Status   string

	Updated string
}

// ExtractImages converts a page of List results into a slice of usable Image structs.
func ExtractImages(page pagination.Page) ([]Image, error) {
	casted := page.(ListPage).Body
	var results []Image

	err := mapstructure.Decode(results, casted)
	return results, err
}
