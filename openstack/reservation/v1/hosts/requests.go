package hosts

import (
	"context"
	"fmt"
	"maps"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToHostListQuery() (string, error)
}

// ListOpts allows the filtering of paginated collections through the API.
// Blazar accepts query parameters on this endpoint but ignores them and
// returns every host.
type ListOpts struct{}

// ToHostListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToHostListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List retrieves a list of compute hosts.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToHostListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return HostPage{pagination.SinglePageBase(r)}
	})
}

// CreateOptsBuilder allows extensions to add parameters to the Create request.
type CreateOptsBuilder interface {
	ToHostCreateMap() (map[string]any, error)
}

// CreateOpts specifies the parameters for enrolling a host into the freepool.
type CreateOpts struct {
	// Name is the name by which Nova knows the hypervisor to enroll.
	Name string `json:"name" required:"true"`

	// ExtraCapabilities holds the operator-defined capabilities to set on the host.
	ExtraCapabilities map[string]any `json:"-"`
}

// ToHostCreateMap formats a CreateOpts into a request body.
func (opts CreateOpts) ToHostCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	for k, v := range opts.ExtraCapabilities {
		if _, ok := b[k]; ok {
			return nil, fmt.Errorf("extra capability %q collides with a host field", k)
		}
		b[k] = v
	}

	return b, nil
}

// Create enrols a compute host into the Blazar freepool.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToHostCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Post(ctx, createURL(client), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Get retrieves a specific compute host based on its unique ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add parameters to the Update request.
type UpdateOptsBuilder interface {
	ToHostUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the parameters for updating a host. Only the extra
// capabilities of a host can be changed.
type UpdateOpts struct {
	// ExtraCapabilities holds the operator-defined capabilities to set on the host.
	ExtraCapabilities map[string]any
}

// ToHostUpdateMap formats an UpdateOpts into a request body.
func (opts UpdateOpts) ToHostUpdateMap() (map[string]any, error) {
	if len(opts.ExtraCapabilities) == 0 {
		return nil, gophercloud.ErrMissingInput{Argument: "ExtraCapabilities"}
	}

	b := make(map[string]any, len(opts.ExtraCapabilities))
	maps.Copy(b, opts.ExtraCapabilities)

	return b, nil
}

// Update changes the extra capabilities of a compute host.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToHostUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Put(ctx, updateURL(client, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete removes a compute host from the Blazar freepool.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListAllocationsOptsBuilder allows extensions to add parameters to the
// ListAllocations request.
type ListAllocationsOptsBuilder interface {
	ToAllocationListQuery() (string, error)
}

// ListAllocationsOpts filters the allocations of the hosts in the freepool.
type ListAllocationsOpts struct {
	// LeaseID only returns the allocations of the given lease.
	LeaseID string `q:"lease_id"`

	// ReservationID only returns the allocations of the given reservation.
	ReservationID string `q:"reservation_id"`
}

// ToAllocationListQuery formats a ListAllocationsOpts into a query string.
func (opts ListAllocationsOpts) ToAllocationListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListAllocations retrieves the reservations holding each compute host, for
// leases that have not ended yet.
func ListAllocations(client *gophercloud.ServiceClient, opts ListAllocationsOptsBuilder) pagination.Pager {
	url := listAllocationsURL(client)
	if opts != nil {
		query, err := opts.ToAllocationListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return AllocationPage{pagination.SinglePageBase(r)}
	})
}

// GetAllocationOptsBuilder allows extensions to add parameters to the
// GetAllocation request.
type GetAllocationOptsBuilder interface {
	ToAllocationGetQuery() (string, error)
}

// GetAllocationOpts filters the allocations of a compute host.
type GetAllocationOpts struct {
	// LeaseID only returns the allocations of the given lease.
	LeaseID string `q:"lease_id"`

	// ReservationID only returns the allocations of the given reservation.
	ReservationID string `q:"reservation_id"`
}

// ToAllocationGetQuery formats a GetAllocationOpts into a query string.
func (opts GetAllocationOpts) ToAllocationGetQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// GetAllocation retrieves the reservations holding a specific compute host,
// for leases that have not ended yet.
func GetAllocation(ctx context.Context, client *gophercloud.ServiceClient, id string, opts GetAllocationOptsBuilder) (r GetAllocationResult) {
	url := getAllocationURL(client, id)
	if opts != nil {
		query, err := opts.ToAllocationGetQuery()
		if err != nil {
			r.Err = err
			return
		}
		url += query
	}

	resp, err := client.Get(ctx, url, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListResourcePropertiesOptsBuilder allows extensions to add parameters to the
// ListResourceProperties request.
type ListResourcePropertiesOptsBuilder interface {
	ToResourcePropertyListQuery() (string, error)
}

// ListResourcePropertiesOpts controls which resource properties are listed.
type ListResourcePropertiesOpts struct {
	// Detail also returns the values each property takes.
	Detail bool `q:"detail"`

	// All also returns the private properties, for an administrator.
	All bool `q:"all"`
}

// ToResourcePropertyListQuery formats a ListResourcePropertiesOpts into a
// query string.
func (opts ListResourcePropertiesOpts) ToResourcePropertyListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// ListResourceProperties retrieves the extra capabilities of the compute hosts
// that reservations can filter on.
func ListResourceProperties(client *gophercloud.ServiceClient, opts ListResourcePropertiesOptsBuilder) pagination.Pager {
	url := listResourcePropertiesURL(client)
	if opts != nil {
		query, err := opts.ToResourcePropertyListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return ResourcePropertyPage{pagination.SinglePageBase(r)}
	})
}

// UpdateResourcePropertyOptsBuilder allows extensions to add parameters to the
// UpdateResourceProperty request.
type UpdateResourcePropertyOptsBuilder interface {
	ToResourcePropertyUpdateMap() (map[string]any, error)
}

// UpdateResourcePropertyOpts specifies the parameters for updating a resource
// property.
type UpdateResourcePropertyOpts struct {
	// Private hides the property from users who list the resource properties.
	Private bool `json:"private"`
}

// ToResourcePropertyUpdateMap formats an UpdateResourcePropertyOpts into a
// request body.
func (opts UpdateResourcePropertyOpts) ToResourcePropertyUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// UpdateResourceProperty changes a resource property of the compute hosts.
func UpdateResourceProperty(ctx context.Context, client *gophercloud.ServiceClient, name string, opts UpdateResourcePropertyOptsBuilder) (r UpdateResourcePropertyResult) {
	b, err := opts.ToResourcePropertyUpdateMap()
	if err != nil {
		r.Err = err
		return
	}

	resp, err := client.Patch(ctx, updateResourcePropertyURL(client, name), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
