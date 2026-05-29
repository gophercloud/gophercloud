package federation

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
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
	ToMappingCreateMap() (map[string]any, error)
}

// UpdateMappingOpts provides options for creating a mapping.
type CreateMappingOpts struct {
	// The list of rules used to map remote users into local users
	Rules []MappingRule `json:"rules"`
}

// ToMappingCreateMap formats a CreateMappingOpts into a create request.
func (opts CreateMappingOpts) ToMappingCreateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "mapping")
}

// CreateMapping creates a new Mapping.
func CreateMapping(ctx context.Context, client *gophercloud.ServiceClient, mappingID string, opts CreateMappingOptsBuilder) (r CreateMappingResult) {
	b, err := opts.ToMappingCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Put(ctx, mappingsResourceURL(client, mappingID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// GetMapping retrieves details on a single mapping, by ID.
func GetMapping(ctx context.Context, client *gophercloud.ServiceClient, mappingID string) (r GetMappingResult) {
	resp, err := client.Get(ctx, mappingsResourceURL(client, mappingID), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateMappingOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateMappingOptsBuilder interface {
	ToMappingUpdateMap() (map[string]any, error)
}

// UpdateMappingOpts provides options for updating a mapping.
type UpdateMappingOpts struct {
	// The list of rules used to map remote users into local users
	Rules []MappingRule `json:"rules"`
}

// ToMappingUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateMappingOpts) ToMappingUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "mapping")
}

// UpdateMapping updates an existing mapping.
func UpdateMapping(ctx context.Context, client *gophercloud.ServiceClient, mappingID string, opts UpdateMappingOptsBuilder) (r UpdateMappingResult) {
	b, err := opts.ToMappingUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := client.Patch(ctx, mappingsResourceURL(client, mappingID), &b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// DeleteMapping deletes a mapping.
func DeleteMapping(ctx context.Context, client *gophercloud.ServiceClient, mappingID string) (r DeleteMappingResult) {
	resp, err := client.Delete(ctx, mappingsResourceURL(client, mappingID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
