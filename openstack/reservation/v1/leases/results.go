package leases

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type Lease struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	UserID    string    `json:"user_id"`
	ProjectID string    `json:"project_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	// The date when this lease was created.
	CreatedAt time.Time `json:"-"`
	// The date when this lease was last updated
	UpdatedAt    time.Time     `json:"-"`
	TrustID      string        `json:"trust_id"`
	Reservations []Reservation `json:"reservations"`
	Events       []Event       `json:"events"`
	Status       string        `json:"status"`
	Degraded     bool          `json:"degraded"`
}

// UnmarshalJSON sets *l to a copy of data.
func (l *Lease) UnmarshalJSON(b []byte) error {
	type tmp Lease
	var s struct {
		tmp
		StartDate gophercloud.JSONRFC3339MilliNoZ `json:"start_date"`
		EndDate   gophercloud.JSONRFC3339MilliNoZ `json:"end_date"`
		CreatedAt gophercloud.JSONRFC3339ZNoTNoZ  `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339ZNoTNoZ  `json:"updated_at"`
	}

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*l = Lease(s.tmp)

	l.StartDate = time.Time(s.StartDate)
	l.EndDate = time.Time(s.EndDate)
	l.CreatedAt = time.Time(s.CreatedAt)
	l.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}

type Event struct {
	ID        string    `json:"id"`
	LeaseID   string    `json:"lease_id"`
	EventType string    `json:"event_type"`
	Time      time.Time `json:"time"`
	// Should change for const Statuses?
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UnmarshalJSON sets *l to a copy of data.
func (e *Event) UnmarshalJSON(b []byte) error {
	type tmp Event
	var s struct {
		tmp
		Time      gophercloud.JSONRFC3339MilliNoZ `json:"time"`
		CreatedAt gophercloud.JSONRFC3339ZNoTNoZ  `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339ZNoTNoZ  `json:"updated_at"`
	}

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*e = Event(s.tmp)

	e.Time = time.Time(s.Time)
	e.CreatedAt = time.Time(s.CreatedAt)
	e.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}

type Reservation struct {
	ID               string    `json:"id"`
	LeaseID          string    `json:"lease_id"`
	ResourceID       string    `json:"resource_id"`
	ResourceType     string    `json:"resource_type"`
	Status           string    `json:"status"`
	MissingResources bool      `json:"missing_resources"`
	ResourcesChanged bool      `json:"resources_changed"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// InstanceReservation fields.
	// InstanceReservation InstanceReservation `json:"instance_reservation"`
	VCPUs              int    `json:"vcpus,omitempty"`
	MemoryMB           int    `json:"memory_mb,omitempty"`
	DiskGB             int    `json:"disk_gb,omitempty"`
	Amount             int    `json:"amount,omitempty"`
	Affinity           bool   `json:"affinity,omitempty"`
	ResourceProperties string `json:"resource_properties,omitempty"`
	FlavorID           string `json:"flavor_id,omitempty"`
	AggregateID        int    `json:"aggregate_id,omitempty"`
	ServerGroupID      string `json:"server_group_id,omitempty"`

	// ComputeHostReservation fields.
	// AggregateID            int                    `json:"aggregate_id"`
	// ResourceProperties     string                 `json:"resource_properties"`
	Min                  int    `json:"min,omitempty"`
	Max                  int    `json:"max,omitempty"`
	HypervisorProperties string `json:"hypervisor_properties,omitempty"`
	BeforeEnd            string `json:"before_end,omitempty"`

	ComputeHostAllocations []string `json:"computehost_allocations"`
	// FloatingIPReservation fields.
	// FloatingIPReservation FloatingIPReservation `json:"floatingip_reservation"`
	NetworkID string `json:"network_id,omitempty"`
	// Amount              int                  `json:"amount"`
	// Required floating IPs addresses.
	RequiredFloatingIPs []string `json:"required_floatingips,omitempty"`

	FloatingIPAllocations []string `json:"floatingip_allocations"`
}

// UnmarshalJSON sets *u to a copy of data.
func (r *Reservation) UnmarshalJSON(b []byte) error {
	type tmp Reservation
	var s struct {
		tmp
		CreatedAt gophercloud.JSONRFC3339ZNoTNoZ `json:"created_at"`
		UpdatedAt gophercloud.JSONRFC3339ZNoTNoZ `json:"updated_at"`
	}

	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*r = Reservation(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}

type leaseResult struct {
	gophercloud.Result
}

type GetResult struct {
	leaseResult
}

type CreateResult struct {
	leaseResult
}

type LeasePage struct {
	pagination.LinkedPageBase
}

func (r LeasePage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	is, err := ExtractLeases(r)
	return len(is) == 0, err
}

func (r leaseResult) Extract() (*Lease, error) {
	var s Lease
	err := r.ExtractInto(&s)
	return &s, err
}

func (r leaseResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "lease")
}

func ExtractLeases(r pagination.Page) ([]Lease, error) {
	var s []Lease
	err := ExtractLeasesInto(r, &s)
	return s, err
}

func ExtractLeasesInto(r pagination.Page, v any) error {
	return r.(LeasePage).Result.ExtractIntoSlicePtr(v, "leases")
}
