package flavors

import (
	"fmt"
	"github.com/racker/perigee"
)

var ErrNotImplemented = fmt.Errorf("Flavors functionality not implemented.")

type ListResults map[string]interface{}
type GetResults map[string]interface{}

// ListFilterOptions helps control the results returned by the List() function.
// ChangesSince, if provided, instructs List to return only those things which have changed since the timestamp provided.
// MinDisk and MinRam, if provided, elides flavors which do not meet your criteria.
// For example, a flavor with a minDisk field of 10 will not be returned if you specify MinDisk set to 20.
// Marker and Limit control paging.
// Limit instructs List to refrain from sending excessively large lists of flavors.
// Marker instructs List where to start listing from.
// Typically, software will use the last ID of the previous call to List to set the Marker for the current call.
type ListFilterOptions struct {
	ChangesSince string
	MinDisk, MinRam int
	Marker string
	Limit int
}

// List instructs OpenStack to provide a list of flavors.
// You may provide criteria by which List curtails its results for easier processing.
// See ListFilterOptions for more details.
func List(c *Client, lfo ListFilterOptions) (ListResults, error) {
	var lr ListResults

	h, err := c.getListHeaders()
	if err != nil {
		return nil, err
	}

	err = perigee.Get(c.getListUrl(lfo), perigee.Options{
		Results:     &lr,
		MoreHeaders: h,
	})
	return lr, err
}

// Get instructs OpenStack to provide details on a single flavor, identified by its ID.
func Get(c *Client, id string) (GetResults, error) {
	var gr GetResults
	h, err := c.getListHeaders()	// same for Get Flavor API
	if err != nil {
		return gr, err
	}
	err = perigee.Get(c.getGetUrl(id), perigee.Options{
		Results: &gr,
		MoreHeaders: h,
	})
	return gr, err
}
