package storagepools

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/gophercloud/gophercloud/pagination"
)

// Minimum set of driver capabilities only
// https://github.com/openstack/cinder/blob/master/cinder/interface/volume_driver.py#L56
type Capabilities struct {
	// Required Fields
	DriverVersion     string  `json:"driver_version"`
	FreeCapacityGB    float64 `json:"-"`
	StorageProtocol   string  `json:"storage_protocol"`
	TotalCapacityGB   float64 `json:"-"`
	VendorName        string  `json:"vendor_name"`
	VolumeBackendName string  `json:"volume_backend_name"`
	// Optional Fields
	ReservedPercentage       int64   `json:"reserved_percentage"`
	LocationInfo             string  `json:"location_info"`
	QoSSupport               bool    `json:"QoS_support"`
	ProvisionedCapacityGB    float64 `json:"provisioned_capacity_gb"`
	MaxOverSubscriptionRatio float64 `json:"max_over_subscription_ratio"`
	ThinProvisioningSupport  bool    `json:"thin_provisioning_support"`
	ThickProvisioningSupport bool    `json:"thick_provisioning_support"`
	TotalVolumes             int64   `json:"total_volumes"`
	FilterFunction           string  `json:"filter_function"`
	GoodnessFuction          string  `json:"goodness_function"`
	Mutliattach              bool    `json:"multiattach"`
	SparseCopyVolume         bool    `json:"sparse_copy_volume"`
}

type StoragePool struct {
	Name         string       `json:"name"`
	Capabilities Capabilities `json:"capabilities"`
}

func (r *StoragePool) UnmarshalJSON(b []byte) error {
	// Unmarshal the generic stuff
	type tmp StoragePool
	var s *struct {
		tmp
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = StoragePool(s.tmp)

	// Unmarshal the more complex things
	var q *struct {
		Capabilities struct {
			FreeCapacityGB  interface{} `json:"free_capacity_gb"`
			TotalCapacityGB interface{} `json:"total_capacity_gb"`
		} `json:"capabilities"`
	}
	err = json.Unmarshal(b, &q)
	if err != nil {
		return err
	}

	// Should be a numeric, "unknown", "infinite"
	if q.Capabilities.FreeCapacityGB != nil {
		switch t := q.Capabilities.FreeCapacityGB.(type) {
		case float64:
			r.Capabilities.FreeCapacityGB = q.Capabilities.FreeCapacityGB.(float64)
		case string:
			keyword := q.Capabilities.FreeCapacityGB.(string)
			switch keyword {
			case "infinite":
				r.Capabilities.FreeCapacityGB = math.Inf(1)
			default:
				r.Capabilities.FreeCapacityGB = 0.0
			}
		default:
			return fmt.Errorf("capabilities.free_capacity_gb: unexpected type %v", t)
		}
	}

	// Should be a numeric, "unknown", "infinite"
	if q.Capabilities.TotalCapacityGB != nil {
		switch t := q.Capabilities.TotalCapacityGB.(type) {
		case float64:
			r.Capabilities.TotalCapacityGB = q.Capabilities.TotalCapacityGB.(float64)
		case string:
			keyword := q.Capabilities.TotalCapacityGB.(string)
			switch keyword {
			case "infinite":
				r.Capabilities.TotalCapacityGB = math.Inf(1)
			default:
				r.Capabilities.TotalCapacityGB = 0.0
			}
		default:
			return fmt.Errorf("capabilities.total_capacity_gb: unexpected type %v", t)
		}
	}

	return nil
}

type StoragePoolPage struct {
	pagination.SinglePageBase
}

func (page StoragePoolPage) IsEmpty() (bool, error) {
	va, err := ExtractStoragePools(page)
	return len(va) == 0, err
}

func ExtractStoragePools(p pagination.Page) ([]StoragePool, error) {
	var s struct {
		StoragePools []StoragePool `json:"pools"`
	}
	err := (p.(StoragePoolPage)).ExtractInto(&s)
	return s.StoragePools, err
}
