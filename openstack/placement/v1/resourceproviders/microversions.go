package resourceproviders

import "github.com/gophercloud/gophercloud/v2"

// ListOpts139 allows filtering resource providers. Filtering is achieved by
// passing in struct field values that map to the resource provider
// attributes you want to see returned.
// ListOpts139 is available in version >= 1.39.
type ListOpts139 struct {
	// Name is the name of the resource provider to filter the list
	Name string `q:"name"`

	// UUID is the uuid of the resource provider to filter the list
	UUID string `q:"uuid"`

	// MemberOf is a list representing aggregate uuids that a provider must be
	// associated with to be returned.
	// Alternative is defined using the in: syntax, e.g. member_of=in:agg1,agg2,agg3.
	// Forbidden aggregates are prefixed with !.
	// Starting with microversion 1.24, the member_of parameter may be repeated.
	MemberOf []string `q:"member_of"`

	// Resources is a comma-separated list of string indicating an amount of resource
	// of a specified class that a provider must have the capacity and availability to serve
	Resources string `q:"resources"`

	// InTree is a string that represents a resource provider UUID.  The returned resource
	// providers will be in the same provider tree as the specified provider.
	InTree string `q:"in_tree"`

	// Required is a list of trait names.
	// Microversion 1.39 added support for repeating the required parameter and for the in: syntax.
	Required []string `q:"required"`
}

// ToResourceProviderListQuery formats a ListOpts139 into a query string.
func (opts ListOpts139) ToResourceProviderListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}
