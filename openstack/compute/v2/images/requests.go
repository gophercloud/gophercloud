package images

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListPage contains a single page of results from a List operation.
// Use ExtractImages to convert it into a slice of usable structs.
type ListPage struct {
	pagination.MarkerPageBase
}

// IsEmpty returns true if a page contains no Image results.
func (page ListPage) IsEmpty() (bool, error) {
	images, err := ExtractImages(page)
	if err != nil {
		return true, err
	}
	return len(images) == 0, nil
}

// LastMarker returns the ID of the final Image on the current page of ListPage.
func (page ListPage) LastMarker() (string, error) {
	images, err := ExtractImages(page)
	if err != nil {
		return "", err
	}
	if len(images) == 0 {
		return "", nil
	}
	return images[len(images)-1].ID, nil
}

// GetResult opaquely stores the result of a Get call.
// Use ExtractImage() to translate it into this provider's version of an Image structure.
type GetResult map[string]interface{}

// List enumerates the available images.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		p := ListPage{pagination.MarkerPageBase{LastHTTPResponse: r}}
		p.MarkerPageBase.Owner = p
		return p
	}

	return pagination.NewPager(client, getListURL(client), createPage)
}

// Get acquires additional detail about a specific image by ID.
// Use ExtractImage() to intepret the result as an openstack Image.
func Get(client *gophercloud.ServiceClient, id string) (GetResult, error) {
	var result GetResult
	_, err := perigee.Request("GET", getImageURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		Results:     &result,
		OkCodes:     []int{200},
	})
	return result, err
}
