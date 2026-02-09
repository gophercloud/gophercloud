package networksegmentranges

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOpts allows filtering when listing network segment ranges.
type ListOpts struct {
	ID              string `q:"id"`
	Name            string `q:"name"`
	Default         *bool  `q:"default"`
	Shared          *bool  `q:"shared"`
	ProjectID       string `q:"project_id"`
	NetworkType     string `q:"network_type"`
	PhysicalNetwork string `q:"physical_network"`
	Minimum         int    `q:"minimum"`
	Maximum         int    `q:"maximum"`
	SortDir         string `q:"sort_dir"`
	SortKey         string `q:"sort_key"`
}

type ListOptsBuilder interface {
	ToNetworkSegmentRangeListQuery() (string, error)
}

func (opts ListOpts) ToNetworkSegmentRangeListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), nil
}

// List retrieves a list of network segment ranges.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToNetworkSegmentRangeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return NetworkSegmentRangePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

type CreateOptsBuilder interface {
	ToNetworkSegmentRangeCreateMap() (map[string]any, error)
}

// CreateOpts represents options for creating a network segment range.
type CreateOpts struct {
	Name            string `json:"name,omitempty"`
	Default         *bool  `json:"default,omitempty"`
	Shared          *bool  `json:"shared,omitempty"`
	ProjectID       string `json:"project_id,omitempty"`
	NetworkType     string `json:"network_type" required:"true"`
	PhysicalNetwork string `json:"physical_network,omitempty"`
	Minimum         int    `json:"minimum" required:"true"`
	Maximum         int    `json:"maximum" required:"true"`
}

func (opts CreateOpts) ToNetworkSegmentRangeCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "network_segment_range")
}

// Create creates a new network segment range.
func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToNetworkSegmentRangeCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, createURL(c), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a specific network segment range.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOptsBuilder interface {
	ToNetworkSegmentRangeUpdateMap() (map[string]any, error)
}

// UpdateOpts represents options for updating a network segment range.
type UpdateOpts struct {
	Name    *string `json:"name,omitempty"`
	Minimum *int    `json:"minimum,omitempty"`
	Maximum *int    `json:"maximum,omitempty"`
}

func (opts UpdateOpts) ToNetworkSegmentRangeUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "network_segment_range")
}

// Update modifies a network segment range.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToNetworkSegmentRangeUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, updateURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete removes a network segment range.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
