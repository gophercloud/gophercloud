package clusters

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// CreateOptsBuilder Builder.
type CreateOptsBuilder interface {
	ToClusterCreateMap() (map[string]interface{}, error)
}

// CreateOpts params
type CreateOpts struct {
	Name            string                 `json:"name" required:"true"`
	DesiredCapacity int                    `json:"desired_capacity"`
	ProfileID       string                 `json:"profile_id" required:"true"`
	MinSize         *int                   `json:"min_size,omitempty"`
	Timeout         int                    `json:"timeout,omitempty"`
	MaxSize         int                    `json:"max_size,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Config          map[string]interface{} `json:"config,omitempty"`
}

// ToClusterCreateMap constructs a request body from CreateOpts.
func (opts CreateOpts) ToClusterCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "cluster")
}

// Create requests the creation of a new cluster.
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToClusterCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	var result *http.Response
	result, r.Err = client.Post(createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})

	if r.Err == nil {
		r.Header = result.Header
	}

	return
}

// Get retrieves details of a single cluster. Use Extract to convert its
// result into a Cluster.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	var result *http.Response
	result, r.Err = client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{OkCodes: []int{200}})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// ListOptsBuilder Builder.
type ListOptsBuilder interface {
	ToClusterListQuery() (string, error)
}

// ListOpts params
type ListOpts struct {
	Limit         int    `q:"limit"`
	Marker        string `q:"marker"`
	Sort          string `q:"sort"`
	GlobalProject *bool  `q:"global_project"`
	Name          string `q:"name,omitempty"`
	Status        string `q:"status,omitempty"`
}

// ToClusterListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToClusterListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List instructs OpenStack to provide a list of clusters.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToClusterListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ClusterPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// UpdateOpts implements UpdateOpts
type UpdateOpts struct {
	Config      string                 `json:"config,omitempty"`
	Name        string                 `json:"name,omitempty"`
	ProfileID   string                 `json:"profile_id,omitempty"`
	Timeout     *int                   `json:"timeout,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	ProfileOnly *bool                  `json:"profile_only,omitempty"`
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToClusterUpdateMap() (map[string]interface{}, error)
}

// ToClusterUpdateMap assembles a request body based on the contents of
// UpdateOpts.
func (opts UpdateOpts) ToClusterUpdateMap() (map[string]interface{}, error) {
	b, err := gophercloud.BuildRequestBody(opts, "cluster")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// Update implements profile updated request.
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToClusterUpdateMap()
	if err != nil {
		r.Err = err
		return r
	}

	var result *http.Response
	result, r.Err = client.Patch(updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 202},
	})

	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// Delete deletes the specified cluster ID.
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	var result *http.Response
	result, r.Err = client.Delete(deleteURL(client, id), nil)
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}

// ResizeOpts params
type ResizeOpts struct {
	AdjustmentType string      `json:"adjustment_type,omitempty"`
	Number         interface{} `json:"number,omitempty"`
	MinSize        *int        `json:"min_size,omitempty"`
	MaxSize        *int        `json:"max_size,omitempty"`
	MinStep        *int        `json:"min_step,omitempty"`
	Strict         *bool       `json:"strict,omitempty"`
}

// ToClusterResizeMap constructs a request body from ResizeOpts.
func (opts ResizeOpts) ToClusterResizeMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "resize")
}

func Resize(client *gophercloud.ServiceClient, id string, opts ResizeOpts) (r ResizeResult) {
	if opts.AdjustmentType != "" && opts.Number == nil {
		r.Err = fmt.Errorf("Number field MUST NOT be empty when AdjustmentType field used")
		return
	}

	switch opts.Number.(type) {
	case int, int32, int64, *int, *int32, *int64:
		// Valid type. Always allow
	case float32, float64, *float32, *float64:
		if opts.AdjustmentType != "CHANGE_IN_PERCENTAGE" {
			r.Err = fmt.Errorf("Only AdjustmentType=CHANGE_IN_PERCENTAGE allows float value for Number field"+
				"AdjustmentType=%s", opts.AdjustmentType)
			return
		}
	default:
		r.Err = fmt.Errorf("Number field must be either int or float %s", reflect.TypeOf(opts.Number))
		return
	}

	b, err := opts.ToClusterResizeMap()
	if err != nil {
		r.Err = err
		return
	}

	var result *http.Response
	result, r.Err = client.Post(actionURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201, 202},
	})
	if r.Err == nil {
		r.Header = result.Header
	}
	return
}
