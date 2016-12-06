package provider

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToProviderNetworkListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the network attributes you want to see returned. SortKey allows you to sort
// by a particular network attribute. SortDir sets the direction, and is either
// `asc' or `desc'. Marker and Limit are used for pagination.
type ListOpts struct {
	Status          string `q:"status"`
	Name            string `q:"name"`
	AdminStateUp    *bool  `q:"admin_state_up"`
	TenantID        string `q:"tenant_id"`
	Shared          *bool  `q:"shared"`
	ID              string `q:"id"`
	NetworkType     string `q:"provider:network_type"`
	PhysicalNetwork string `q:"provider:physical_network"`
	SegmentationID  string `q:"provider:segmentation_id"`
	Marker          string `q:"marker"`
	Limit           int    `q:"limit"`
	SortKey         string `q:"sort_key"`
	SortDir         string `q:"sort_dir"`
}

// ToProviderNetworkListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToProviderNetworkListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// provider networks. It accepts a ListOpts struct, which allows you to filter and sort
// the returned collection for greater efficiency.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToProviderNetworkListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return NetworkPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific provider network based on its unique ID.
func Get(c *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = c.Get(getURL(c, id), &r.Body, nil)
	return
}

// CreateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Create operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type CreateOptsBuilder interface {
	ToProviderNetworkCreateMap() (map[string]interface{}, error)
}

// CreateOpts satisfies the CreateOptsBuilder interface
type CreateOpts struct {
	AdminStateUp    *bool  `json:"admin_state_up,omitempty"`
	Name            string `json:"name,omitempty"`
	Shared          *bool  `json:"shared,omitempty"`
	TenantID        string `json:"tenant_id,omitempty"`
	NetworkType     string `json:"provider:network_type,omitempty"`
	PhysicalNetwork string `json:"provider:physical_network,omitempty"`
	SegmentationID  string `json:"provider:segmentation_id,omitempty"`
}

// ToProviderNetworkCreateMap casts a CreateOpts struct to a map.
func (opts CreateOpts) ToProviderNetworkCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "network")
}
// Create accepts a CreateOpts struct and creates a new provider network using the values
// provided. This operation does not actually require a request body, i.e. the
// CreateOpts struct argument can be empty.
//
// The tenant ID that is contained in the URI is the tenant that creates the
// network. An admin user, however, has the option of specifying another tenant
// ID in the CreateOpts struct.
func Create(c *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToProviderNetworkCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}

// UpdateOptsBuilder is the interface options structs have to satisfy in order
// to be used in the main Update operation in this package. Since many
// extensions decorate or modify the common logic, it is useful for them to
// satisfy a basic interface in order for them to be used.
type UpdateOptsBuilder interface {
	ToProviderNetworkUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts satisfies the UpdateOptsBuilder interface
type UpdateOpts struct {
	AdminStateUp    *bool  `json:"admin_state_up,omitempty"`
	Name            string `json:"name,omitempty"`
	Shared          *bool  `json:"shared,omitempty"`
	NetworkType     string `json:"provider:network_type"`
	PhysicalNetwork string `json:"provider:physical_network"`
	SegmentationID  string `json:"provider:segmentation_id"`
}

// ToProviderNetworkUpdateMap casts a UpdateOpts struct to a map.
func (opts UpdateOpts) ToProviderNetworkUpdateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "network")
}

// Update accepts a UpdateOpts struct and updates an existing provider network using the
// values provided. For more information, see the Create function.
func Update(c *gophercloud.ServiceClient, networkID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToProviderNetworkUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(updateURL(c, networkID), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200, 201},
	})
	return
}

// Delete accepts a unique ID and deletes the provider network associated with it.
func Delete(c *gophercloud.ServiceClient, networkID string) (r DeleteResult) {
	_, r.Err = c.Delete(deleteURL(c, networkID), nil)
	return
}

// IDFromName is a convenience function that returns a provider network's ID given its name.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""
	pages, err := List(client, nil).AllPages()
	if err != nil {
		return "", err
	}

	all, err := ExtractNetworkExtAttrs(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gophercloud.ErrResourceNotFound{Name: name, ResourceType: "network"}
	case 1:
		return id, nil
	default:
		return "", gophercloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "network"}
	}
}
