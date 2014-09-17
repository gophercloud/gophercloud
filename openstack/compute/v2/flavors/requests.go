package flavors

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// ListPage contains a single page of the response from a List call.
type ListPage struct {
	pagination.MarkerPageBase
}

// IsEmpty determines if a page contains any results.
func (p ListPage) IsEmpty() (bool, error) {
	flavors, err := ExtractFlavors(p)
	if err != nil {
		return true, err
	}
	return len(flavors) == 0, nil
}

// LastMarker returns the ID field of the final result from this page, to be used as the marker for the next.
func (p ListPage) LastMarker() (string, error) {
	flavors, err := ExtractFlavors(p)
	if err != nil {
		return "", err
	}
	if len(flavors) == 0 {
		return "", nil
	}
	return flavors[len(flavors)-1].ID, nil
}

// GetResults temporarily encodes the result of a Get operation.
type GetResults map[string]interface{}

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
		p := ListPage{pagination.MarkerPageBase{LastHTTPResponse: r}}
		p.MarkerPageBase.Owner = p
		return p
	}

	return pagination.NewPager(client, getListURL(client, lfo), createPage)
}

// Get instructs OpenStack to provide details on a single flavor, identified by its ID.
// Use ExtractFlavor to convert its result into a Flavor.
func Get(client *gophercloud.ServiceClient, id string) (GetResults, error) {
	var gr GetResults
	err := perigee.Get(getFlavorURL(client, id), perigee.Options{
		Results:     &gr,
		MoreHeaders: client.Provider.AuthenticatedHeaders(),
	})
	return gr, err
}
