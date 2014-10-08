package images

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListOpts contain options for limiting the number of Images returned from a call to ListDetail.
type ListOpts struct {
	// When the image last changed status (in date-time format).
	ChangesSince string `q:"changes-since"`
	// The number of Images to return.
	Limit int `q:"limit"`
	// UUID of the Image at which to set a marker.
	Marker string `q:"marker"`
	// The name of the Image.
	Name string `q:"name:"`
	// The name of the Server (in URL format).
	Server string `q:"server"`
	// The current status of the Image.
	Status string `q:"status"`
	// The value of the type of image (e.g. BASE, SERVER, ALL)
	Type string `q:"type"`
}

// ListDetail enumerates the available images.
func ListDetail(client *gophercloud.ServiceClient, opts *ListOpts) pagination.Pager {
	url := listDetailURL(client)
	if opts != nil {
		query, err := gophercloud.BuildQueryString(opts)
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query.String()
	}

	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ImagePage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	}

	return pagination.NewPager(client, url, createPage)
}

// Get acquires additional detail about a specific image by ID.
// Use ExtractImage() to intepret the result as an openstack Image.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var result GetResult
	_, result.Err = perigee.Request("GET", getURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		Results:     &result.Resp,
		OkCodes:     []int{200},
	})
	return result
}
