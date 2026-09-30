package leases

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// dateFormat is the layout Blazar requires for the dates in a lease request.
const dateFormat = "2006-01-02 15:04"

// The resource types Blazar is able to reserve.
const (
	// ResourceTypeHost reserves whole compute hosts from the freepool.
	ResourceTypeHost = "physical:host"

	// ResourceTypeInstance reserves capacity for instances of a size given
	// in the reservation.
	ResourceTypeInstance = "virtual:instance"

	// ResourceTypeFlavorInstance reserves capacity for instances of an
	// existing Nova flavor.
	ResourceTypeFlavorInstance = "flavor:instance"
)

// ListOptsBuilder allows extensions to add parameters to the List request.
type ListOptsBuilder interface {
	ToLeaseListQuery() (string, error)
}

// ListOpts allows the filtering of paginated collections through the API.
// Blazar accepts query parameters on this endpoint but ignores them and
// returns every lease the caller can see.
type ListOpts struct{}

// ToLeaseListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToLeaseListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List retrieves a list of leases. An administrator sees every lease, any
// other caller sees only the leases of their own project.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToLeaseListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return LeasePage{pagination.SinglePageBase(r)}
	})
}

// ReservationOptsBuilder allows a reservation of any resource type to be
// added to a lease.
type ReservationOptsBuilder interface {
	ToReservationMap() (map[string]any, error)
}

// HostReservationOpts reserves whole compute hosts from the Blazar freepool.
type HostReservationOpts struct {
	// Min is the smallest number of hosts the lease can be satisfied with.
	// It must be at least 1.
	Min int `json:"min" required:"true"`

	// Max is the largest number of hosts to reserve. It must be at least Min.
	Max int `json:"max" required:"true"`

	// HypervisorProperties filters the candidate hosts on the properties Nova
	// reports, e.g. `[">=", "$vcpus", "4"]`. Empty matches every host.
	HypervisorProperties string `json:"hypervisor_properties"`

	// ResourceProperties filters the candidate hosts on their extra
	// capabilities, in the same form as HypervisorProperties.
	ResourceProperties string `json:"resource_properties"`

	// BeforeEnd is the action taken when the before_end_lease event fires:
	// "snapshot", "default" or "" for no action. Nil means "default".
	BeforeEnd *string `json:"before_end,omitempty"`
}

// ToReservationMap formats a HostReservationOpts into a reservation body.
func (opts HostReservationOpts) ToReservationMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	b["resource_type"] = ResourceTypeHost

	return b, nil
}

// InstanceReservationOpts reserves capacity for instances of the given size.
type InstanceReservationOpts struct {
	// Amount is the number of instances to reserve capacity for.
	Amount int `json:"amount" required:"true"`

	// VCPUs is the number of virtual CPUs per instance.
	VCPUs int `json:"vcpus" required:"true"`

	// MemoryMB is the amount of memory per instance, in megabytes.
	MemoryMB int `json:"memory_mb" required:"true"`

	// DiskGB is the amount of disk per instance, in gigabytes.
	DiskGB int `json:"disk_gb" required:"true"`

	// Affinity places the instances on the same host (true) or on distinct
	// hosts (false). Nil leaves the placement to Nova.
	Affinity *bool `json:"affinity"`

	// ResourceProperties filters the candidate hosts on their extra
	// capabilities, e.g. `["==", "$gpu", "a100"]`. Empty matches every host.
	ResourceProperties string `json:"resource_properties"`
}

// ToReservationMap formats an InstanceReservationOpts into a reservation body.
func (opts InstanceReservationOpts) ToReservationMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	b["resource_type"] = ResourceTypeInstance

	return b, nil
}

// FlavorInstanceReservationOpts reserves capacity for instances of an
// existing Nova flavor.
type FlavorInstanceReservationOpts struct {
	// Amount is the number of instances to reserve capacity for.
	Amount int `json:"amount" required:"true"`

	// FlavorID is the ID of the Nova flavor to reserve instances of. Blazar
	// derives the size and the resource properties of the reservation from it.
	FlavorID string `json:"flavor_id" required:"true"`
}

// ToReservationMap formats a FlavorInstanceReservationOpts into a reservation
// body.
func (opts FlavorInstanceReservationOpts) ToReservationMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	// Blazar reads affinity for a flavor reservation but rejects every value
	// other than null, so there is nothing to expose on the opts.
	b["affinity"] = nil
	b["resource_type"] = ResourceTypeFlavorInstance

	return b, nil
}

// CreateOptsBuilder allows extensions to add parameters to the Create request.
type CreateOptsBuilder interface {
	ToLeaseCreateMap() (map[string]any, error)
}

// CreateOpts specifies the parameters for creating a lease.
type CreateOpts struct {
	// Name is the name of the lease. It must be unique within the project.
	Name string `json:"name" required:"true"`

	// StartDate is the time at which the lease becomes active. Blazar
	// truncates it to the minute and rejects a time in the past.
	StartDate time.Time `json:"-"`

	// EndDate is the time at which the lease expires. It must be later than
	// StartDate.
	EndDate time.Time `json:"-"`

	// BeforeEndDate is the time at which the before_end_lease event fires.
	// It must fall within the lease. Blazar picks a default when it is nil.
	BeforeEndDate *time.Time `json:"-"`

	// Reservations holds the resources to reserve for the duration of the
	// lease. At least one is required.
	Reservations []ReservationOptsBuilder `json:"-"`
}

// ToLeaseCreateMap formats a CreateOpts into a request body.
func (opts CreateOpts) ToLeaseCreateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.StartDate.IsZero() {
		return nil, gophercloud.ErrMissingInput{Argument: "StartDate"}
	}
	if opts.EndDate.IsZero() {
		return nil, gophercloud.ErrMissingInput{Argument: "EndDate"}
	}
	if len(opts.Reservations) == 0 {
		return nil, gophercloud.ErrMissingInput{Argument: "Reservations"}
	}

	b["start_date"] = opts.StartDate.UTC().Format(dateFormat)
	b["end_date"] = opts.EndDate.UTC().Format(dateFormat)
	if opts.BeforeEndDate != nil {
		b["before_end_date"] = opts.BeforeEndDate.UTC().Format(dateFormat)
	}

	reservations := make([]map[string]any, len(opts.Reservations))
	for i, reservation := range opts.Reservations {
		if reservation == nil {
			return nil, gophercloud.ErrMissingInput{Argument: fmt.Sprintf("Reservations[%d]", i)}
		}

		reservations[i], err = reservation.ToReservationMap()
		if err != nil {
			return nil, err
		}
	}
	b["reservations"] = reservations

	return b, nil
}

// Create creates a lease and the reservations it holds.
func Create(ctx context.Context, client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToLeaseCreateMap()
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

// Get retrieves a specific lease based on its unique ID.
func Get(ctx context.Context, client *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := client.Get(ctx, getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateReservationOpts specifies the changes to make to one reservation of a
// lease. Only the fields to change need to be set.
type UpdateReservationOpts struct {
	// ID is the unique identifier of the reservation to update. It must
	// belong to the lease being updated.
	ID string `json:"id" required:"true"`

	// Min is the smallest number of hosts the lease can be satisfied with,
	// for a host reservation.
	Min int `json:"min,omitempty"`

	// Max is the largest number of hosts to reserve, for a host reservation.
	Max int `json:"max,omitempty"`

	// HypervisorProperties replaces the Nova property filter of a host
	// reservation. A pointer to the empty string clears it.
	HypervisorProperties *string `json:"hypervisor_properties,omitempty"`

	// ResourceProperties replaces the extra capability filter of a host or
	// instance reservation. A pointer to the empty string clears it.
	ResourceProperties *string `json:"resource_properties,omitempty"`

	// Amount is the number of instances to reserve capacity for, for an
	// instance reservation.
	Amount int `json:"amount,omitempty"`

	// VCPUs is the number of virtual CPUs per instance, for an instance
	// reservation.
	VCPUs int `json:"vcpus,omitempty"`

	// MemoryMB is the amount of memory per instance in megabytes, for an
	// instance reservation.
	MemoryMB int `json:"memory_mb,omitempty"`

	// DiskGB is the amount of disk per instance in gigabytes, for an instance
	// reservation.
	DiskGB int `json:"disk_gb,omitempty"`

	// Affinity replaces the placement constraint of an instance reservation.
	Affinity *bool `json:"affinity,omitempty"`
}

// UpdateOptsBuilder allows extensions to add parameters to the Update request.
type UpdateOptsBuilder interface {
	ToLeaseUpdateMap() (map[string]any, error)
}

// UpdateOpts specifies the parameters for updating a lease. Every field is
// optional, and Blazar keeps the current value of anything left unset.
type UpdateOpts struct {
	// Name is the new name of the lease. Renaming is the only change allowed
	// on a lease that has already ended.
	Name string `json:"name,omitempty"`

	// StartDate is the new time at which the lease becomes active. Blazar
	// rejects it once the lease has started.
	StartDate *time.Time `json:"-"`

	// EndDate is the new time at which the lease expires.
	EndDate *time.Time `json:"-"`

	// BeforeEndDate is the new time at which Blazar fires the
	// before_end_lease event. It must fall within the lease.
	BeforeEndDate *time.Time `json:"-"`

	// Reservations holds the changes to make to the reservations of the lease.
	Reservations []UpdateReservationOpts `json:"reservations,omitempty"`
}

// ToLeaseUpdateMap formats an UpdateOpts into a request body.
func (opts UpdateOpts) ToLeaseUpdateMap() (map[string]any, error) {
	b, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	if opts.StartDate != nil {
		b["start_date"] = opts.StartDate.UTC().Format(dateFormat)
	}
	if opts.EndDate != nil {
		b["end_date"] = opts.EndDate.UTC().Format(dateFormat)
	}
	if opts.BeforeEndDate != nil {
		b["before_end_date"] = opts.BeforeEndDate.UTC().Format(dateFormat)
	}

	if len(b) == 0 {
		return nil, gophercloud.ErrMissingInput{Argument: "UpdateOpts"}
	}

	return b, nil
}

// Update changes a lease and the reservations it holds.
func Update(ctx context.Context, client *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToLeaseUpdateMap()
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

// Delete removes a lease and releases the resources it reserved.
func Delete(ctx context.Context, client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := client.Delete(ctx, deleteURL(client, id), &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
