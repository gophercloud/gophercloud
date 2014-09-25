package images

import (
	"github.com/mitchellh/mapstructure"
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListPage contains a single page of results from a List operation.
// Use ExtractImages to convert it into a slice of usable structs.
type ListPage struct {
	pagination.LinkedPageBase
}

// IsEmpty returns true if a page contains no Image results.
func (page ListPage) IsEmpty() (bool, error) {
	images, err := ExtractImages(page)
	if err != nil {
		return true, err
	}
	return len(images) == 0, nil
}

// NextPageURL uses the response's embedded link reference to navigate to the next page of results.
func (page ListPage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"images_links"`
	}

	var r resp
	err := mapstructure.Decode(page.Body, &r)
	if err != nil {
		return "", err
	}

	var url string
	for _, l := range r.Links {
		if l.Rel == "next" {
			url = l.Href
		}
	}
	if url == "" {
		return "", nil
	}

	return url, nil
}

// List enumerates the available images.
func List(client *gophercloud.ServiceClient) pagination.Pager {
	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ListPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	}

	return pagination.NewPager(client, getListURL(client), createPage)
}

// Get acquires additional detail about a specific image by ID.
// Use ExtractImage() to intepret the result as an openstack Image.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var result GetResult
	_, result.Err = perigee.Request("GET", getImageURL(client, id), perigee.Options{
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
		Results:     &result.Resp,
		OkCodes:     []int{200},
	})
	return result
}
