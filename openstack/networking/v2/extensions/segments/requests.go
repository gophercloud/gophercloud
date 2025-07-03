package segments

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOpts allows filtering when listing segments.
type ListOpts struct {
	Name            string `q:"name"`
	Description     string `q:"description"`
	NetworkID       string `q:"network_id"`
	PhysicalNetwork string `q:"physical_network"`
	NetworkType     string `q:"network_type"`
	SegmentationID  int    `q:"segmentation_id"`
	RevisionNumber  int    `q:"revision_number"`
	SortDir         string `q:"sort_dir"`
	SortKey         string `q:"sort_key"`
	Fields          string `q:"fields"`
}

// ListOptsBuilder interface for listing.
type ListOptsBuilder interface {
	ToSegmentListQuery() (string, error)
}

func (opts ListOpts) ToSegmentListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List all segments.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(c)
	if opts != nil {
		query, err := opts.ToSegmentListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return SegmentPage{LinkedPageBase: pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateOptsBuilder allows extensions to add additional parameters.
type CreateOptsBuilder interface {
	ToSegmentCreateMap() (map[string]any, error)
}

// CreateOpts contains the fields needed for creating a segment.
type CreateOpts struct {
	Name            string `json:"name,omitempty"`
	Description     string `json:"description,omitempty"`
	NetworkID       string `json:"network_id" required:"true"`
	NetworkType     string `json:"network_type" required:"true"`
	PhysicalNetwork string `json:"physical_network,omitempty"`
	SegmentationID  int    `json:"segmentation_id,omitempty"`
}

func (opts CreateOpts) ToSegmentCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "segment")
}

// Create a new segment.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToSegmentCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, rootURL(c), b, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a segment by ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, resourceURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete removes a segment by ID.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, resourceURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOpts contains fields to update a segment.
type UpdateOpts struct {
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	SegmentationID *int    `json:"segmentation_id,omitempty"`
}

// UpdateOptsBuilder is the interface for update options.
type UpdateOptsBuilder interface {
	ToSegmentUpdateMap() (map[string]any, error)
}

func (opts UpdateOpts) ToSegmentUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "segment")
}

// Update a segment.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToSegmentUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, resourceURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
