package allocationcandidates

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToAllocationCandidatesListQuery() (string, error)
}

// ResourceGroup represents a granular resource request group.
// It is not a standalone request type, but becomes part of ListOpts,
// keyed by a user-defined suffix string.
// Available in version >= 1.25.
type ResourceGroup struct {
	// Resources is a comma-separated list of resource amounts, e.g. "VCPU:1,MEMORY_MB:1024".
	// Becomes the resourcesN query parameter.
	Resources string

	// Required is a list of trait expressions for this group.
	// Each entry becomes one requiredN query parameter.
	// Available in version >= 1.39: can be repeated; supports in: syntax.
	Required []string

	// MemberOf is an aggregate UUID or the prefix in: followed by a
	// comma-separated list of aggregate UUIDs for this group.
	// Becomes the member_ofN query parameter.
	// Available in version >= 1.32: forbidden aggregates can be expressed
	// with a ! prefix or the !in: prefix.
	MemberOf string

	// InTree filters results to include only providers in the same tree as
	// the specified provider UUID for this group.
	// Becomes the in_treeN query parameter.
	// Available in version >= 1.31.
	InTree string
}

// ListOpts allows filtering of allocation candidates.
type ListOpts struct {
	// Resources is a comma-separated list of resource amounts that providers
	// must collectively have capacity to serve, e.g. "VCPU:4,DISK_GB:64,MEMORY_MB:2048".
	Resources string `q:"resources"`

	// Required is a comma-separated list of traits that a provider must have.
	// Available in version >= 1.17.
	// Available in version >= 1.22: prefix with ! for forbidden traits.
	// Available in version >= 1.39: can be repeated and supports in: syntax.
	Required []string `q:"required"`

	// MemberOf is a string representing an aggregate UUID, or the prefix in:
	// followed by a comma-separated list of aggregate UUIDs.
	// Available in version >= 1.21.
	// Available in version >= 1.24: can be specified multiple times.
	MemberOf []string `q:"member_of"`

	// InTree is a resource provider UUID. When supplied, filters candidates to
	// only those providers that are in the same tree.
	// Available in version >= 1.31.
	InTree string `q:"in_tree"`

	// GroupPolicy indicates how the groups should interact when more than one
	// resourcesN parameter is supplied. Valid values are "none" and "isolate".
	// Available in version >= 1.25.
	GroupPolicy string `q:"group_policy"`

	// Limit is a positive integer used to limit the maximum number of
	// allocation candidates returned.
	// Available in version >= 1.16.
	Limit int `q:"limit"`

	// RootRequired is a comma-separated list of trait requirements that the
	// root provider of the (non-sharing) tree must satisfy.
	// Available in version >= 1.35.
	RootRequired string `q:"root_required"`

	// SameSubtree is a comma-separated list of request group suffix strings.
	// At least one of the resource providers satisfying a specified request group
	// must be an ancestor of the rest.
	// Available in version >= 1.36.
	SameSubtree []string `q:"same_subtree"`

	// ResourceGroups allows specifying suffixed granular resource request groups.
	// The map key is the non-empty group suffix (e.g. "1", "_NET1", "_STORAGE").
	// Use the top-level Resources/Required/MemberOf/InTree fields for the
	// unsuffixed (default) group; do not use an empty string key here.
	// In microversions 1.25-1.32 the suffix must be a numeric string.
	// Starting from microversion 1.33 it can be 1-64 characters [a-zA-Z0-9_-].
	// Available in version >= 1.25.
	ResourceGroups map[string]ResourceGroup `q:"-"`
}

// ToAllocationCandidatesListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToAllocationCandidatesListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	if err != nil {
		return "", err
	}

	params := q.Query()
	for suffix, group := range opts.ResourceGroups {
		if group.Resources != "" {
			params.Add("resources"+suffix, group.Resources)
		}
		for _, required := range group.Required {
			if required != "" {
				params.Add("required"+suffix, required)
			}
		}
		if group.MemberOf != "" {
			params.Add("member_of"+suffix, group.MemberOf)
		}
		if group.InTree != "" {
			params.Add("in_tree"+suffix, group.InTree)
		}
	}
	q.RawQuery = params.Encode()

	return q.String(), nil
}

// List makes a request against the API to list allocation candidates.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)

	if opts != nil {
		query, err := opts.ToAllocationCandidatesListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AllocationCandidatesPage{pagination.SinglePageBase(r)}
	})
}
