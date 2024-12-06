package stackresources

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Find retrieves stack resources for the given stack name.
func Find(ctx context.Context, c gophercloud.Client, stackName string) (r FindResult) {
	resp, err := c.Get(ctx, findURL(c, stackName), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToStackResourceListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Marker and Limit are used for pagination.
type ListOpts struct {
	// Include resources from nest stacks up to Depth levels of recursion.
	Depth int `q:"nested_depth"`
}

// ToStackResourceListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToStackResourceListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list resources for the given stack.
func List(client gophercloud.Client, stackName, stackID string, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client, stackName, stackID)
	if opts != nil {
		query, err := opts.ToStackResourceListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ResourcePage{pagination.SinglePageBase(r)}
	})
}

// Get retreives data for the given stack resource.
func Get(ctx context.Context, c gophercloud.Client, stackName, stackID, resourceName string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, stackName, stackID, resourceName), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Metadata retreives the metadata for the given stack resource.
func Metadata(ctx context.Context, c gophercloud.Client, stackName, stackID, resourceName string) (r MetadataResult) {
	resp, err := c.Get(ctx, metadataURL(c, stackName, stackID, resourceName), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListTypes makes a request against the API to list resource types.
func ListTypes(client gophercloud.Client) pagination.Pager {
	return pagination.NewPager(client, listTypesURL(client), func(r pagination.PageResult) pagination.Page {
		return ResourceTypePage{pagination.SinglePageBase(r)}
	})
}

// Schema retreives the schema for the given resource type.
func Schema(ctx context.Context, c gophercloud.Client, resourceType string) (r SchemaResult) {
	resp, err := c.Get(ctx, schemaURL(c, resourceType), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Template retreives the template representation for the given resource type.
func Template(ctx context.Context, c gophercloud.Client, resourceType string) (r TemplateResult) {
	resp, err := c.Get(ctx, templateURL(c, resourceType), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// MarkUnhealthyOpts contains the common options struct used in this package's
// MarkUnhealthy operations.
type MarkUnhealthyOpts struct {
	// A boolean indicating whether the target resource should be marked as unhealthy.
	MarkUnhealthy bool `json:"mark_unhealthy"`
	// The reason for the current stack resource state.
	ResourceStatusReason string `json:"resource_status_reason,omitempty"`
}

// MarkUnhealthyOptsBuilder is the interface options structs have to satisfy in order
// to be used in the MarkUnhealthy operation in this package
type MarkUnhealthyOptsBuilder interface {
	ToMarkUnhealthyMap() (map[string]any, error)
}

// ToMarkUnhealthyMap validates that a template was supplied and calls
// the ToMarkUnhealthyMap private function.
func (opts MarkUnhealthyOpts) ToMarkUnhealthyMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}
	return b, nil
}

// MarkUnhealthy marks the specified resource in the stack as unhealthy.
func MarkUnhealthy(ctx context.Context, c gophercloud.Client, stackName, stackID, resourceName string, opts MarkUnhealthyOptsBuilder) (r MarkUnhealthyResult) {
	b, err := opts.ToMarkUnhealthyMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Patch(ctx, markUnhealthyURL(c, stackName, stackID, resourceName), b, nil, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
