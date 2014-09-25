package flavors

import (
	"github.com/mitchellh/mapstructure"
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListPage contains a single page of the response from a List call.
type ListPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines if a page contains any results.
func (p ListPage) IsEmpty() (bool, error) {
	flavors, err := ExtractFlavors(p)
	if err != nil {
		return true, err
	}
	return len(flavors) == 0, nil
}

// NextPageURL uses the response's embedded link reference to navigate to the next page of results.
func (p ListPage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"flavors_links"`
	}

	var r resp
	err := mapstructure.Decode(p.Body, &r)
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

// ListFilterOptions helps control the results returned by the List() function.
// For example, a flavor with a minDisk field of 10 will not be returned if you specify MinDisk set to 20.
// Typically, software will use the last ID of the previous call to List to set the Marker for the current call.
type ListFilterOptions struct {

	// ChangesSince, if provided, instructs List to return only those things which have changed since the timestamp provided.
	ChangesSince string

	// MinDisk and MinRAM, if provided, elides flavors which do not meet your criteria.
	MinDisk, MinRAM int

	// Marker and Limit control paging.
	// Marker instructs List where to start listing from.
	Marker string

	// Limit instructs List to refrain from sending excessively large lists of flavors.
	Limit int
}

// List instructs OpenStack to provide a list of flavors.
// You may provide criteria by which List curtails its results for easier processing.
// See ListFilterOptions for more details.
func List(client *gophercloud.ServiceClient, lfo ListFilterOptions) pagination.Pager {
	createPage := func(r pagination.LastHTTPResponse) pagination.Page {
		return ListPage{pagination.LinkedPageBase{LastHTTPResponse: r}}
	}

	return pagination.NewPager(client, getListURL(client, lfo), createPage)
}

// Get instructs OpenStack to provide details on a single flavor, identified by its ID.
// Use ExtractFlavor to convert its result into a Flavor.
func Get(client *gophercloud.ServiceClient, id string) GetResult {
	var gr GetResult
	gr.Err = perigee.Get(getFlavorURL(client, id), perigee.Options{
		Results:     &gr.Resp,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	return gr
}
