package schedulerstats

import (
	"encoding/json"
	"math"

	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Capabilities represents the information of an individual Pool.
type Capabilities struct {
	// The following fields should be present in all storage drivers.

	// The quality of service (QoS) support.
	Qos bool `json:"qos"`
	// The date and time stamp when the API request was issued.
	Timestamp string `json:"timestamp"`
	// The name of the share back end.
	ShareBackendName string `json:"share_backend_name"`
	// Share server is usually a storage virtual machine or a lightweight container that is used to export shared file systems.
	DriverHandlesShareServers bool `json:"driver_handles_share_servers"`
	// The driver version of the back end.
	DriverVersion string `json:"driver_version"`
	// The amount of free capacity for the back end, in GiBs. A valid value is a string, such as unknown, or an integer.
	FreeCapacityGB float64 `json:"-"`
	// The storage protocol for the back end. For example, NFS_CIFS, glusterfs, HDFS, etc.
	StorageProtocol string `json:"storage_protocol"`
	// The total capacity for the back end, in GiBs. A valid value is a string, such as unknown, or an integer.
	TotalCapacityGB float64 `json:"-"`
	// The specification that filters back ends by whether they do or do not support share snapshots.
	SnapshotSupport bool `json:"snapshot_support"`
	// The back end replication domain.
	ReplicationDomain string `json:"replication_domain"`
	// The name of the vendor for the back end.
	VendorName string `json:"vendor_name"`

	// The following fields are optional and may have empty values depending

	// on the storage driver in use.
	ReservedPercentage  int64   `json:"reserved_percentage"`
	AllocatedCapacityGB float64 `json:"-"`
}

// Pool represents an individual Pool retrieved from the
// schedulerstats API.
type Pool struct {
	// The name of the back end.
	Name string `json:"name"`
	// The name of the back end.
	Backend string `json:"backend"`
	// The pool name for the back end.
	Pool string `json:"pool"`
	// The host name for the back end.
	Host string `json:"host"`
	// The back end capabilities which include qos, total_capacity_gb, etc.
	Capabilities Capabilities `json:"capabilities,omitempty"`
}

func (r *Capabilities) UnmarshalJSON(b []byte) error {
	type tmp Capabilities
	var s struct {
		tmp
		AllocatedCapacityGB any `json:"allocated_capacity_gb"`
		FreeCapacityGB      any `json:"free_capacity_gb"`
		TotalCapacityGB     any `json:"total_capacity_gb"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Capabilities(s.tmp)

	// Generic function to parse a capacity value which may be a numeric
	// value, "unknown", or "infinite"
	parseCapacity := func(capacity any) float64 {
		if capacity != nil {
			switch c := capacity.(type) {
			case float64:
				return c
			case string:
				if c == "infinite" {
					return math.Inf(1)
				}
			}
		}
		return 0.0
	}

	r.AllocatedCapacityGB = parseCapacity(s.AllocatedCapacityGB)
	r.FreeCapacityGB = parseCapacity(s.FreeCapacityGB)
	r.TotalCapacityGB = parseCapacity(s.TotalCapacityGB)

	return nil
}

// PoolPage is a single page of all List results.
type PoolPage struct {
	pagination.SinglePageBase
}

// IsEmpty satisfies the IsEmpty method of the Page interface. It returns true
// if a List contains no results.
func (page PoolPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}

	va, err := ExtractPools(page)
	return len(va) == 0, err
}

// ExtractPools takes a List result and extracts the collection of
// Pools returned by the API.
func ExtractPools(p pagination.Page) ([]Pool, error) {
	var s struct {
		Pools []Pool `json:"pools"`
	}
	err := (p.(PoolPage)).ExtractInto(&s)
	return s.Pools, err
}
