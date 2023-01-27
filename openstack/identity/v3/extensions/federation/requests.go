package federation

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListMappings enumerates the mappings.
func ListMappings(client *gophercloud.ServiceClient) pagination.Pager {
	return pagination.NewPager(client, mappingsRootURL(client), func(r pagination.PageResult) pagination.Page {
		return MappingsPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// CreateMappingOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateMappingOptsBuilder interface {
	ToMappingCreateMap() (map[string]interface{}, error)
}

// UpdateMappingOpts provides options for creating a mapping.
type CreateMappingOpts struct {
	// The list of rules used to map remote users into local users
	Rules []MappingRule `json:"rules"`
}

// ToMappingCreateMap formats a CreateMappingOpts into a create request.
func (opts CreateMappingOpts) ToMappingCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "mapping")
}

// CreateMapping creates a new Mapping.
func CreateMapping(client *gophercloud.ServiceClient, mappingID string, opts CreateMappingOptsBuilder) (r CreateMappingResult) {
	b, err := opts.ToMappingCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(mappingsResourceURL(client, mappingID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
