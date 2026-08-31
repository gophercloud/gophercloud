package hosts

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

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
