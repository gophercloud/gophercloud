package images

import (
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListResults contains a single page of results from a List operation.
// Use ExtractImages to convert it into a slice of usable structs.
type ListResults struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a page contains no Image results.
func (page ListResults) IsEmpty() (bool, error) {
	images, err := ExtractImages(page)
	if err != nil {
		return true, err
	}
	return len(images) == 0, nil
}

// LastMarker returns the ID of the final Image on the current page of ListResults.
func (page ListResults) LastMarker() (string, error) {
	images, err := ExtractImages(page)
	if err != nil {
		return "", err
	}
	if len(images) == 0 {
		return "", nil
	}
	return images[len(images)-1].ID, nil
}

// List enumerates the available images.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		p := ListResults{pagination.MarkerPageBase{LastHTTPResponse: r}}
		p.MarkerPageBase.Owner = p
		return p
	}

	return pagination.NewPager(client, getListURL(client), createPage)
}
