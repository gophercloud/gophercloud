package schedulerstats

import (
	"encoding/json"
	"math"
	"strconv"

	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Capabilities represents the information of an individual StoragePool.
type Capabilities struct {
	// The following fields should be present in all storage drivers.
	DriverVersion     string  `json:"driver_version"`
	FreeCapacityGB    float64 `json:"-"`
	StorageProtocol   string  `json:"storage_protocol"`
	TotalCapacityGB   float64 `json:"-"`
	VendorName        string  `json:"vendor_name"`
	VolumeBackendName string  `json:"volume_backend_name"`

	// The following fields are optional and may have empty values depending
	// on the storage driver in use.
	ReservedPercentage       int64   `json:"reserved_percentage"`
	LocationInfo             string  `json:"location_info"`
	QoSSupport               bool    `json:"QoS_support"`
	ProvisionedCapacityGB    float64 `json:"provisioned_capacity_gb"`
	MaxOverSubscriptionRatio string  `json:"-"`
	ThinProvisioningSupport  bool    `json:"thin_provisioning_support"`
	ThickProvisioningSupport bool    `json:"thick_provisioning_support"`
	TotalVolumes             int64   `json:"total_volumes"`
	FilterFunction           string  `json:"filter_function"`
	GoodnessFunction         string  `json:"goodness_function"`
	Multiattach              bool    `json:"multiattach"`
	SparseCopyVolume         bool    `json:"sparse_copy_volume"`
	AllocatedCapacityGB      float64 `json:"-"`
}

// StoragePool represents an individual StoragePool retrieved from the
// schedulerstats API.
type StoragePool struct {
	Name         string       `json:"name"`
	Capabilities Capabilities `json:"capabilities"`
}

func (r *Capabilities) UnmarshalJSON(b []byte) error {
	type tmp Capabilities
	var s struct {
		tmp
		AllocatedCapacityGB      any `json:"allocated_capacity_gb"`
		FreeCapacityGB           any `json:"free_capacity_gb"`
		MaxOverSubscriptionRatio any `json:"max_over_subscription_ratio"`
		TotalCapacityGB          any `json:"total_capacity_gb"`
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

	if s.MaxOverSubscriptionRatio != nil {
		switch t := s.MaxOverSubscriptionRatio.(type) {
		case float64:
			r.MaxOverSubscriptionRatio = strconv.FormatFloat(t, 'f', -1, 64)
		case string:
			r.MaxOverSubscriptionRatio = t
		}
	}

	return nil
}

// StoragePoolPage is a single page of all List results.
type StoragePoolPage struct {
	pagination.SinglePageBase
}

// IsEmpty satisfies the IsEmpty method of the Page interface. It returns true
// if a List contains no results.
func (page StoragePoolPage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}

	va, err := ExtractStoragePools(page)
	return len(va) == 0, err
}

// ExtractStoragePools takes a List result and extracts the collection of
// StoragePools returned by the API.
func ExtractStoragePools(p pagination.Page) ([]StoragePool, error) {
	var s struct {
		StoragePools []StoragePool `json:"pools"`
	}
	err := (p.(StoragePoolPage)).ExtractInto(&s)
	return s.StoragePools, err
}
