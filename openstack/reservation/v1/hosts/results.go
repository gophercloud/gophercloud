package hosts

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract interprets any commonResult as a Host.
func (r commonResult) Extract() (*Host, error) {
	var s struct {
		Host *Host `json:"host"`
	}
	err := r.ExtractInto(&s)
	return s.Host, err
}

// GetResult is the response from a Get operation. Call its Extract method to
// interpret it as a Host.
type GetResult struct {
	commonResult
}

// CreateResult is the response from a Create operation. Call its Extract
// method to interpret it as a Host.
type CreateResult struct {
	commonResult
}

// UpdateResult is the response from an Update operation. Call its Extract
// method to interpret it as a Host.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr
// method to determine whether the request succeeded.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Host represents a compute host enrolled in the Blazar freepool.
type Host struct {
	// ID is the unique identifier of the host within Blazar. It is distinct
	// from the Nova hypervisor ID.
	ID string `json:"id"`

	// HypervisorHostname is the name by which Nova knows the hypervisor.
	HypervisorHostname string `json:"hypervisor_hostname"`

	// HypervisorType is the virtualisation technology of the hypervisor.
	HypervisorType string `json:"hypervisor_type"`

	// HypervisorVersion is the version of the hypervisor software.
	HypervisorVersion int `json:"hypervisor_version"`

	// ServiceName is the name of the nova-compute service running on the host.
	ServiceName string `json:"service_name"`

	// VCPUs is the number of virtual CPUs the host provides.
	VCPUs int `json:"vcpus"`

	// CPUInfo is the hypervisor's description of the host CPU, as a JSON string.
	CPUInfo string `json:"cpu_info"`

	// MemoryMB is the amount of memory the host provides, in megabytes.
	MemoryMB int `json:"memory_mb"`

	// LocalGB is the amount of local disk the host provides, in gigabytes.
	LocalGB int `json:"local_gb"`

	// Status is the current status of the host.
	Status string `json:"status"`

	// AvailabilityZone is the availability zone the host belongs to.
	AvailabilityZone string `json:"availability_zone"`

	// TrustID is the identifier of the Keystone trust Blazar uses to act on
	// the host. It is created by Blazar.
	TrustID string `json:"trust_id"`

	// Reservable reports whether the host may be allocated to new reservations.
	Reservable bool `json:"reservable"`

	// CreatedAt is the time at which the host was enrolled.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the time at which the host was last modified. It is nil
	// if the host has never been modified.
	UpdatedAt *time.Time `json:"-"`

	// ExtraCapabilities holds the operator-defined capabilities of the host,
	// which Blazar flattens into the top level of the host object.
	ExtraCapabilities map[string]any `json:"-"`
}

// UnmarshalJSON implements unmarshalling custom types
func (r *Host) UnmarshalJSON(b []byte) error {
	type tmp Host
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339ZNoTNoZ  `json:"created_at"`
		UpdatedAt *gophercloud.JSONRFC3339ZNoTNoZ `json:"updated_at"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*r = Host(s.tmp)
	r.CreatedAt = time.Time(s.CreatedAt)
	if s.UpdatedAt != nil {
		t := time.Time(*s.UpdatedAt)
		r.UpdatedAt = &t
	}

	var resultMap map[string]any
	if err := json.Unmarshal(b, &resultMap); err != nil {
		return err
	}

	delete(resultMap, "created_at")
	delete(resultMap, "updated_at")
	r.ExtraCapabilities = gophercloud.RemainingKeys(Host{}, resultMap)

	return nil
}

// HostPage contains a single page of all hosts from a List call.
type HostPage struct {
	pagination.SinglePageBase
}

func (r HostPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	hosts, err := ExtractHosts(r)
	return len(hosts) == 0, err
}

// ExtractHosts takes a List result and extracts the collection of hosts
// returned by the API.
func ExtractHosts(p pagination.Page) ([]Host, error) {
	var s struct {
		Hosts []Host `json:"hosts"`
	}
	err := (p.(HostPage)).ExtractInto(&s)
	return s.Hosts, err
}

// GetAllocationResult is the response from a GetAllocation operation. Call its
// Extract method to interpret it as an Allocation.
type GetAllocationResult struct {
	gophercloud.Result
}

// Extract interprets a GetAllocationResult as an Allocation.
func (r GetAllocationResult) Extract() (*Allocation, error) {
	var s struct {
		Allocation *Allocation `json:"allocation"`
	}
	err := r.ExtractInto(&s)
	return s.Allocation, err
}

// Allocation represents the reservations holding a compute host.
type Allocation struct {
	// ResourceID is the unique identifier of the host within Blazar.
	ResourceID string `json:"resource_id"`

	// Reservations holds the reservations of the leases that have not ended
	// yet and hold the host.
	Reservations []AllocationReservation `json:"reservations"`
}

// AllocationReservation represents a reservation holding a compute host.
type AllocationReservation struct {
	// ID is the unique identifier of the reservation.
	ID string `json:"id"`

	// LeaseID is the identifier of the lease holding the reservation.
	LeaseID string `json:"lease_id"`

	// StartDate is the time at which the lease becomes active.
	StartDate time.Time `json:"-"`

	// EndDate is the time at which the lease expires.
	EndDate time.Time `json:"-"`
}

// UnmarshalJSON implements unmarshalling custom types
func (r *AllocationReservation) UnmarshalJSON(b []byte) error {
	type tmp AllocationReservation
	var s struct {
		tmp
		StartDate gophercloud.JSONRFC3339MilliNoZ `json:"start_date"`
		EndDate   gophercloud.JSONRFC3339MilliNoZ `json:"end_date"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*r = AllocationReservation(s.tmp)
	r.StartDate = time.Time(s.StartDate)
	r.EndDate = time.Time(s.EndDate)

	return nil
}

// AllocationPage contains a single page of all allocations from a
// ListAllocations call.
type AllocationPage struct {
	pagination.SinglePageBase
}

func (r AllocationPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	allocations, err := ExtractAllocations(r)
	return len(allocations) == 0, err
}

// ExtractAllocations takes a ListAllocations result and extracts the
// collection of allocations returned by the API.
func ExtractAllocations(p pagination.Page) ([]Allocation, error) {
	var s struct {
		Allocations []Allocation `json:"allocations"`
	}
	err := (p.(AllocationPage)).ExtractInto(&s)
	return s.Allocations, err
}

// ResourceProperty represents an extra capability of the compute hosts that
// reservations can filter on.
type ResourceProperty struct {
	// Property is the name of the extra capability.
	Property string `json:"property"`

	// Private reports whether the property is hidden from users. It is only
	// returned when listing with Detail. Blazar currently reports every
	// property as public, see https://bugs.launchpad.net/blazar/+bug/2169110.
	Private *bool `json:"private"`

	// Values holds the values the hosts have for the property. It is only
	// returned when listing with Detail.
	Values []string `json:"values"`
}

// ResourcePropertyPage contains a single page of all resource properties from
// a ListResourceProperties call.
type ResourcePropertyPage struct {
	pagination.SinglePageBase
}

func (r ResourcePropertyPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	properties, err := ExtractResourceProperties(r)
	return len(properties) == 0, err
}

// ExtractResourceProperties takes a ListResourceProperties result and extracts
// the collection of resource properties returned by the API.
func ExtractResourceProperties(p pagination.Page) ([]ResourceProperty, error) {
	var s struct {
		ResourceProperties []ResourceProperty `json:"resource_properties"`
	}
	err := (p.(ResourcePropertyPage)).ExtractInto(&s)
	return s.ResourceProperties, err
}

// UpdateResourcePropertyResult is the response from an UpdateResourceProperty
// operation. Call its Extract method to interpret it as an
// UpdatedResourceProperty.
type UpdateResourcePropertyResult struct {
	gophercloud.Result
}

// Extract interprets an UpdateResourcePropertyResult as an
// UpdatedResourceProperty.
func (r UpdateResourcePropertyResult) Extract() (*UpdatedResourceProperty, error) {
	var s struct {
		ResourceProperty *UpdatedResourceProperty `json:"resource_property"`
	}
	err := r.ExtractInto(&s)
	return s.ResourceProperty, err
}

// UpdatedResourceProperty represents a resource property as stored by Blazar.
type UpdatedResourceProperty struct {
	// ID is the unique identifier of the resource property.
	ID string `json:"id"`

	// ResourceType is the type of resource the property belongs to.
	ResourceType string `json:"resource_type"`

	// PropertyName is the name of the extra capability.
	PropertyName string `json:"property_name"`

	// Private reports whether the property is hidden from users.
	Private bool `json:"private"`

	// CreatedAt is the time at which the property was first seen.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the time at which the property was last modified. It is nil
	// if the property has never been modified.
	UpdatedAt *time.Time `json:"-"`
}

// UnmarshalJSON implements unmarshalling custom types
func (r *UpdatedResourceProperty) UnmarshalJSON(b []byte) error {
	type tmp UpdatedResourceProperty
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339MilliNoZ  `json:"created_at"`
		UpdatedAt *gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*r = UpdatedResourceProperty(s.tmp)
	r.CreatedAt = time.Time(s.CreatedAt)
	if s.UpdatedAt != nil {
		t := time.Time(*s.UpdatedAt)
		r.UpdatedAt = &t
	}

	return nil
}
