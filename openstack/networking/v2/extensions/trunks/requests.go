package trunks

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToTrunkCreateMap() (map[string]any, error)
}

// CreateOpts represents the attributes used when creating a new trunk.
type CreateOpts struct {
	TenantID     string    `json:"tenant_id,omitempty"`
	ProjectID    string    `json:"project_id,omitempty"`
	PortID       string    `json:"port_id" required:"true"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	AdminStateUp *bool     `json:"admin_state_up,omitempty"`
	Subports     []Subport `json:"sub_ports"`
}

// ToTrunkCreateMap builds a request body from CreateOpts.
func (opts CreateOpts) ToTrunkCreateMap() (map[string]any, error) {
	if opts.Subports == nil {
		opts.Subports = []Subport{}
	}
	return gophercloud.BuildRequestBody(opts, "trunk")
}

func Create(ctx context.Context, c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	body, err := opts.ToTrunkCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := c.Post(ctx, createURL(c), body, &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete accepts a unique ID and deletes the trunk associated with it.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToTrunkListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the trunk attributes you want to see returned. SortKey allows you to sort
// by a particular trunk attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	AdminStateUp *bool  `q:"admin_state_up"`
	Description  string `q:"description"`
	ID           string `q:"id"`
	Name         string `q:"name"`
	PortID       string `q:"port_id"`
	Status       string `q:"status"`
	TenantID     string `q:"tenant_id"`
	ProjectID    string `q:"project_id"`
	SortDir      string `q:"sort_dir"`
	SortKey      string `q:"sort_key"`
	Tags         string `q:"tags"`
	TagsAny      string `q:"tags-any"`
	NotTags      string `q:"not-tags"`
	NotTagsAny   string `q:"not-tags-any"`
	// TODO change type to *int for consistency
	RevisionNumber string `q:"revision_number"`
}

// ToTrunkListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToTrunkListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// trunks. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
//
// Default policy settings return only those trunks that are owned by the tenant
// who submits the request, unless the request is submitted by a user with
// administrative rights.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToTrunkListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return TrunkPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific trunk based on its unique ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type UpdateOptsBuilder interface {
	ToTrunkUpdateMap() (map[string]any, error)
}

type UpdateOpts struct {
	AdminStateUp *bool   `json:"admin_state_up,omitempty"`
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`

	// RevisionNumber implements extension:standard-attr-revisions. If != "" it
	// will set revision_number=%s. If the revision number does not match, the
	// update will fail.
	RevisionNumber *int `json:"-" h:"If-Match"`
}

func (opts UpdateOpts) ToTrunkUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "trunk")
}

func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	body, err := opts.ToTrunkUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	h, err := gophercloud.BuildHeaders(opts)
	if err != nil {
		r.Err = err
		return
	}
	for k := range h {
		if k == "If-Match" {
			h[k] = fmt.Sprintf("revision_number=%s", h[k])
		}
	}
	resp, err := c.Put(ctx, updateURL(c, id), body, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: h,
		OkCodes:     []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

func GetSubports(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetSubportsResult) {
	resp, err := c.Get(ctx, getSubportsURL(c, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type AddSubportsOpts struct {
	Subports []Subport `json:"sub_ports" required:"true"`
}

type AddSubportsOptsBuilder interface {
	ToTrunkAddSubportsMap() (map[string]any, error)
}

func (opts AddSubportsOpts) ToTrunkAddSubportsMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func AddSubports(ctx context.Context, c *gophercloud.ServiceClient, id string, opts AddSubportsOptsBuilder) (r UpdateSubportsResult) {
	body, err := opts.ToTrunkAddSubportsMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, addSubportsURL(c, id), body, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

type RemoveSubport struct {
	PortID string `json:"port_id" required:"true"`
}

type RemoveSubportsOpts struct {
	Subports []RemoveSubport `json:"sub_ports"`
}

type RemoveSubportsOptsBuilder interface {
	ToTrunkRemoveSubportsMap() (map[string]any, error)
}

func (opts RemoveSubportsOpts) ToTrunkRemoveSubportsMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

func RemoveSubports(ctx context.Context, c *gophercloud.ServiceClient, id string, opts RemoveSubportsOptsBuilder) (r UpdateSubportsResult) {
	body, err := opts.ToTrunkRemoveSubportsMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, removeSubportsURL(c, id), body, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
