package leases

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract interprets any commonResult as a Lease.
func (r commonResult) Extract() (*Lease, error) {
	var s struct {
		Lease *Lease `json:"lease"`
	}
	err := r.ExtractInto(&s)
	return s.Lease, err
}

// GetResult is the response from a Get operation. Call its Extract method to
// interpret it as a Lease.
type GetResult struct {
	commonResult
}

// CreateResult is the response from a Create operation. Call its Extract
// method to interpret it as a Lease.
type CreateResult struct {
	commonResult
}

// UpdateResult is the response from an Update operation. Call its Extract
// method to interpret it as a Lease.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the response from a Delete operation. Call its ExtractErr
// method to determine whether the request succeeded.
type DeleteResult struct {
	gophercloud.ErrResult
}

// Lease represents a reservation of resources over a period of time.
type Lease struct {
	// ID is the unique identifier of the lease.
	ID string `json:"id"`

	// Name is the name of the lease, unique within the project.
	Name string `json:"name"`

	// StartDate is the time at which the lease becomes active.
	StartDate time.Time `json:"-"`

	// EndDate is the time at which the lease expires.
	EndDate time.Time `json:"-"`

	// Status is the current status of the lease.
	Status string `json:"status"`

	// Degraded reports whether the lease has lost any of the resources it
	// reserved.
	Degraded bool `json:"degraded"`

	// UserID is the identifier of the user who created the lease.
	UserID string `json:"user_id"`

	// ProjectID is the identifier of the project that owns the lease.
	ProjectID string `json:"project_id"`

	// TrustID is the identifier of the Keystone trust Blazar uses to act on
	// behalf of the lease owner. It is created by Blazar.
	TrustID string `json:"trust_id"`

	// CreatedAt is the time at which the lease was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the time at which the lease was last modified. It is nil
	// if the lease has never been modified.
	UpdatedAt *time.Time `json:"-"`

	// Reservations holds the resource reservations of the lease. Blazar does
	// not order them, so match on ResourceType rather than on position.
	Reservations []Reservation `json:"reservations"`

	// Events holds the state transitions Blazar has scheduled for the lease.
	// They are not ordered, so match on EventType rather than on position.
	Events []Event `json:"events"`
}

// UnmarshalJSON implements unmarshalling custom types
func (r *Lease) UnmarshalJSON(b []byte) error {
	type tmp Lease
	var s struct {
		tmp
		StartDate gophercloud.JSONRFC3339MilliNoZ `json:"start_date"`
		EndDate   gophercloud.JSONRFC3339MilliNoZ `json:"end_date"`
		CreatedAt gophercloud.JSONRFC3339ZNoTNoZ  `json:"created_at"`
		UpdatedAt *gophercloud.JSONRFC3339ZNoTNoZ `json:"updated_at"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*r = Lease(s.tmp)
	r.StartDate = time.Time(s.StartDate)
	r.EndDate = time.Time(s.EndDate)
	r.CreatedAt = time.Time(s.CreatedAt)
	if s.UpdatedAt != nil {
		t := time.Time(*s.UpdatedAt)
		r.UpdatedAt = &t
	}

	return nil
}

// Reservation represents the resources of one resource type reserved by a
// lease. The fields Blazar returns depend on ResourceType, so the ones that
// are specific to a resource type are nil for the others.
type Reservation struct {
	// ID is the unique identifier of the reservation.
	ID string `json:"id"`

	// LeaseID is the identifier of the lease holding the reservation.
	LeaseID string `json:"lease_id"`

	// Status is the current status of the reservation.
	Status string `json:"status"`

	// MissingResources reports whether Blazar was unable to reserve every
	// resource the reservation asked for.
	MissingResources bool `json:"missing_resources"`

	// ResourcesChanged reports whether Blazar had to swap out any of the
	// reserved resources, for instance after a host failed.
	ResourcesChanged bool `json:"resources_changed"`

	// ResourceID is the identifier of the resource Blazar created to back the
	// reservation. It is a host reservation for ResourceTypeHost and an
	// instance reservation otherwise.
	ResourceID string `json:"resource_id"`

	// ResourceType is the type of resource reserved. It is one of
	// ResourceTypeHost, ResourceTypeInstance or ResourceTypeFlavorInstance.
	ResourceType string `json:"resource_type"`

	// Min is the smallest number of hosts the reservation can be satisfied
	// with. It is only returned for a host reservation.
	Min *int `json:"min"`

	// Max is the largest number of hosts reserved. It is only returned for a
	// host reservation.
	Max *int `json:"max"`

	// HypervisorProperties is the expression filtering the candidate hosts on
	// the properties Nova reports. It is only returned for a host reservation.
	HypervisorProperties *string `json:"hypervisor_properties"`

	// ResourceProperties is the expression filtering the candidate hosts on
	// their extra capabilities.
	ResourceProperties *string `json:"resource_properties"`

	// BeforeEnd is the action Blazar takes when the before_end_lease event
	// fires. It is only returned for a host reservation.
	BeforeEnd *string `json:"before_end"`

	// Amount is the number of resources reserved. It is returned for an
	// instance reservation and for a floating IP reservation.
	Amount *int `json:"amount"`

	// VCPUs is the number of virtual CPUs per reserved instance. It is only
	// returned for an instance reservation.
	VCPUs *int `json:"vcpus"`

	// MemoryMB is the amount of memory per reserved instance, in megabytes.
	// It is only returned for an instance reservation.
	MemoryMB *int `json:"memory_mb"`

	// DiskGB is the amount of disk per reserved instance, in gigabytes. It is
	// only returned for an instance reservation.
	DiskGB *int `json:"disk_gb"`

	// Affinity is the placement constraint on the reserved instances. It is
	// nil when the placement is left to Nova.
	Affinity *bool `json:"affinity"`

	// FlavorID is the identifier of the Nova flavor Blazar created for the
	// reserved instances. It is the same as ID.
	FlavorID *string `json:"flavor_id"`

	// ServerGroupID is the identifier of the Nova server group Blazar created
	// to enforce Affinity. It is nil when Affinity is nil.
	ServerGroupID *string `json:"server_group_id"`

	// AggregateID is the identifier of the Nova host aggregate holding the
	// reserved hosts. It is only returned for an instance reservation.
	AggregateID *int `json:"aggregate_id"`

	// CreatedAt is the time at which the reservation was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the time at which the reservation was last modified. It is
	// nil if the reservation has never been modified.
	UpdatedAt *time.Time `json:"-"`
}

// UnmarshalJSON implements unmarshalling custom types
func (r *Reservation) UnmarshalJSON(b []byte) error {
	type tmp Reservation
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339ZNoTNoZ  `json:"created_at"`
		UpdatedAt *gophercloud.JSONRFC3339ZNoTNoZ `json:"updated_at"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*r = Reservation(s.tmp)
	r.CreatedAt = time.Time(s.CreatedAt)
	if s.UpdatedAt != nil {
		t := time.Time(*s.UpdatedAt)
		r.UpdatedAt = &t
	}

	return nil
}

// Event represents a state transition Blazar has scheduled for a lease.
type Event struct {
	// ID is the unique identifier of the event.
	ID string `json:"id"`

	// LeaseID is the identifier of the lease the event belongs to.
	LeaseID string `json:"lease_id"`

	// Status is the current status of the event.
	Status string `json:"status"`

	// EventType is the transition the event performs, such as "start_lease",
	// "end_lease" or "before_end_lease".
	EventType string `json:"event_type"`

	// Time is the time at which the event fires.
	Time time.Time `json:"-"`

	// CreatedAt is the time at which the event was created.
	CreatedAt time.Time `json:"-"`

	// UpdatedAt is the time at which the event was last modified. It is nil
	// if the event has never been modified.
	UpdatedAt *time.Time `json:"-"`
}

// UnmarshalJSON implements unmarshalling custom types
func (r *Event) UnmarshalJSON(b []byte) error {
	type tmp Event
	var s struct {
		tmp
		Time      gophercloud.JSONRFC3339MilliNoZ `json:"time"`
		CreatedAt gophercloud.JSONRFC3339ZNoTNoZ  `json:"created_at"`
		UpdatedAt *gophercloud.JSONRFC3339ZNoTNoZ `json:"updated_at"`
	}
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	*r = Event(s.tmp)
	r.Time = time.Time(s.Time)
	r.CreatedAt = time.Time(s.CreatedAt)
	if s.UpdatedAt != nil {
		t := time.Time(*s.UpdatedAt)
		r.UpdatedAt = &t
	}

	return nil
}

// LeasePage contains a single page of all leases from a List call.
type LeasePage struct {
	pagination.SinglePageBase
}

func (r LeasePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	leases, err := ExtractLeases(r)
	return len(leases) == 0, err
}

// ExtractLeases takes a List result and extracts the collection of leases
// returned by the API.
func ExtractLeases(p pagination.Page) ([]Lease, error) {
	var s struct {
		Leases []Lease `json:"leases"`
	}
	err := (p.(LeasePage)).ExtractInto(&s)
	return s.Leases, err
}
