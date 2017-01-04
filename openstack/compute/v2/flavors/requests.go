package flavors

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToFlavorListQuery() (string, error)
}

// FlavorType maps to OpenStack's Flavor.is_public field. Although the is_public field is boolean, the
// request options are ternary, which is why FlavorType is a string. The following values are
// allowed:
//
//      Project (the default): Returns only public flavors and private flavors associated with that project.
//      Private (admin only):  Returns only private flavors, across all projects.
//      All (admin only):      Returns all public and private flavors across all projects.
//
// If no IsPublic argument is supplied in the request, OpenStack treats the FlavorType as Project.
type FlavorType string
const (
        Project FlavorType =  "true"
        Private FlavorType = "false"
        All FlavorType = "None"
)
func (flavorType FlavorType) String() string {
    switch flavorType {
    case Project:
        return "Project"
    case Private:
        return "Private"
    case All:
        return "All"
    }
    return "Unknown"
}

// ListOpts helps control the results returned by the List() function.
// For example, a flavor with a minDisk field of 10 will not be returned if you specify MinDisk set to 20.
// Typically, software will use the last ID of the previous call to List to set the Marker for the current call.
type ListOpts struct {

	// ChangesSince, if provided, instructs List to return only those things which have changed since the timestamp provided.
	ChangesSince string `q:"changes-since"`

	// MinDisk and MinRAM, if provided, elides flavors which do not meet your criteria.
	MinDisk int `q:"minDisk"`
	MinRAM  int `q:"minRam"`

	// Marker and Limit control paging.
	// Marker instructs List where to start listing from.
	Marker string `q:"marker"`

	// Limit instructs List to refrain from sending excessively large lists of flavors.
	Limit int `q:"limit"`

        // IsPublic, if provided, instructs List which set of flavors to return. If IsPublic not provided,
        // flavors for the current project are returned.
        FlavorType FlavorType `q:"is_public"`
}

// ToFlavorListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToFlavorListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListDetail instructs OpenStack to provide a list of flavors.
// You may provide criteria by which List curtails its results for easier processing.
// See ListOpts for more details.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToFlavorListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return FlavorPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get instructs OpenStack to provide details on a single flavor, identified by its ID.
// Use ExtractFlavor to convert its result into a Flavor.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// IDFromName is a convienience function that returns a flavor's ID given its name.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""
	allPages, err := ListDetail(client, nil).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractFlavors(allPages)
	if err != nil {
		return "", err
	}

	for _, f := range all {
		if f.Name == name {
			count++
			id = f.ID
		}
	}

	switch count {
	case 0:
		err := &gophercloud.ErrResourceNotFound{}
		err.ResourceType = "flavor"
		err.Name = name
		return "", err
	case 1:
		return id, nil
	default:
		err := &gophercloud.ErrMultipleResourcesFound{}
		err.ResourceType = "flavor"
		err.Name = name
		err.Count = count
		return "", err
	}
}
