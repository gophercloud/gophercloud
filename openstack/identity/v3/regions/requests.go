package regions

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToRegionListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// ParentRegionID filters the response by a parent region ID.
	ParentRegionID string `q:"parent_region_id"`
}

// ToRegionListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRegionListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the Regions to which the current token has access.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToRegionListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RegionPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single region, by ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToRegionCreateMap() (map[string]any, error)
}

// CreateOpts provides options used to create a region.
type CreateOpts struct {
	// ID is the ID of the new region.
	ID string `json:"id,omitempty"`

	// Description is a description of the region.
	Description string `json:"description,omitempty"`

	// ParentRegionID is the ID of the parent the region to add this region under.
	ParentRegionID string `json:"parent_region_id,omitempty"`

	// Extra is free-form extra key/value pairs to describe the region.
	Extra map[string]any `json:"-"`
}

// ToRegionCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToRegionCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "region")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["region"].(map[string]any); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create creates a new Region.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRegionCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Post(ctx, createURL(client), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToRegionUpdateMap() (map[string]any, error)
}

// UpdateOpts provides options for updating a region.
type UpdateOpts struct {
	// Description is a description of the region.
	Description *string `json:"description,omitempty"`

	// ParentRegionID is the ID of the parent region.
	ParentRegionID string `json:"parent_region_id,omitempty"`

	/*
		// Due to a bug in Keystone, the Extra column of the Region table
		// is not updatable, see: https://bugs.launchpad.net/keystone/+bug/1729933
		// The following lines should be uncommented once the fix is merged.

		// Extra is free-form extra key/value pairs to describe the region.
		Extra map[string]any `json:"-"`
	*/
}

// ToRegionUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToRegionUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "region")
	if err != nil {
		return nil, err
	}

	/*
		// Due to a bug in Keystone, the Extra column of the Region table
		// is not updatable, see: https://bugs.launchpad.net/keystone/+bug/1729933
		// The following lines should be uncommented once the fix is merged.

		if opts.Extra != nil {
			if v, ok := b["region"].(map[string]any); ok {
				for key, value := range opts.Extra {
					v[key] = value
				}
			}
		}
	*/

	return b, nil
}

// Update updates an existing Region.
func Update(ctx context.Context, client *gophercloud.ServiceClient, regionID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRegionUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, updateURL(client, regionID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a region.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, regionID string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, regionID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
