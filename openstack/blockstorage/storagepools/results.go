package storagepools

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/gophercloud/gophercloud"
)

type ListResult struct {
	gophercloud.Result
}

// Minimum iset of driver capabilities only
// https://github.com/openstack/cinder/blob/master/doc/source/devref/drivers.rst#volume-stats
type Capabilities struct {
	DriverVersion     string  `json:"driver_version"`
	FreeCapacityGB    float64 `json:"-"`
	StorageProtocol   string  `json:"storage_protocol"`
	TotalCapacityGB   float64 `json:"-"`
	VendorName        string  `json:"vendor_name"`
	VolumeBackendName string  `json:"volume_backend_name"`
}

type StoragePool struct {
	Name         string `json:"name"`
	Capabilities `json:"capabilities"`
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
		case int:
			r.Capabilities.FreeCapacityGB = q.Capabilities.FreeCapacityGB.(float64)
		case string:
			keyword := q.Capabilities.FreeCapacityGB.(string)
			switch keyword {
			case "unknown":
				r.Capabilities.FreeCapacityGB = 0.0
			case "infinite":
				r.Capabilities.FreeCapacityGB = math.Inf(1)
			default:
				return fmt.Errorf("capabilities.free_capacity_gb: unexpected string %v", keyword)
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
		case int:
			r.Capabilities.TotalCapacityGB = q.Capabilities.TotalCapacityGB.(float64)
		case string:
			keyword := q.Capabilities.TotalCapacityGB.(string)
			switch keyword {
			case "unknown":
				r.Capabilities.TotalCapacityGB = 0.0
			case "infinite":
				r.Capabilities.TotalCapacityGB = math.Inf(1)
			default:
				return fmt.Errorf("capabilities.free_capacity_gb: unexpected string %v", keyword)
			}
		default:
			return fmt.Errorf("capabilities.total_capacity_gb: unexpected type %v", t)
		}
	}

	return nil
}

func ExtractStoragePools(r ListResult) ([]StoragePool, error) {
	var s struct {
		StoragePools []StoragePool `json:"pools"`
	}
	err := r.ExtractInto(&s)
	return s.StoragePools, err
}
