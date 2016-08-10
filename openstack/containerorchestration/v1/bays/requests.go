package bays

import (
	"encoding/json"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/containerorchestration/v1/common"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToBayListQuery() (string, error)
}

// ListOpts allows the sorting of paginated collections through
// the API. SortKey allows you to sort by a particular bay attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for pagination.
type ListOpts struct {
	Marker  string `q:"marker"`
	Limit   int    `q:"limit"`
	SortKey string `q:"sort_key"`
	SortDir string `q:"sort_dir"`
}

// ToBayListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToBayListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// bays. It accepts a ListOptsBuilder, which allows you to sort
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToBayListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return BayPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific bay based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToBayCreateMap() (map[string]interface{}, error)
}

// CreateOpts satisfies the CreateOptsBuilder interface
type CreateOpts struct {
	BayModelID   string `json:"baymodel_id" required:"true"`
	Name         string `json:"name,omitempty"`
	Masters      int    `json:"master_count,omitempty"`
	Nodes        int    `json:"node_count,omitempty"`
	DiscoveryURL string `json:"discovery_url,omitempty"`
}

// ToBayCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToBayCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// Create accepts a CreateOpts struct and creates a new bay using the values
// provided.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToBayCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

type ErrDeleteFailed struct { }

// Error400 extracts the actual error message from the body of the response
func (d ErrDeleteFailed) Error400(r gophercloud.ErrUnexpectedResponseCode) error {
	var s *common.ErrorResponse
	err := json.Unmarshal(r.Body, &s)
	if err != nil {
		return gophercloud.ErrDefault400{ErrUnexpectedResponseCode: r}
	}

	return s
}

func (d ErrDeleteFailed) Error() string {
	return "Unable to delete bay"
}

// Delete accepts a unique ID and deletes the bay associated with it.
func Delete(c *gophercloud.ServiceClient, bayID string) (r DeleteResult) {
	opts := &gophercloud.RequestOpts{ErrorContext: ErrDeleteFailed{}}
	_, r.Err = c.Delete(deleteURL(c, bayID), opts)
	return
}
