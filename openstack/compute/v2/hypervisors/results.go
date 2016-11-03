package hypervisors

import (
	"encoding/json"
	"fmt"
	"github.com/gophercloud/gophercloud/pagination"
)

type Topology struct {
	Sockets int `json:"sockets"`
	Cores   int `json:"cores"`
	Threads int `json:"threads"`
}

type CpuInfo struct {
	Vendor   string   `json:"vendor"`
	Arch     string   `json:"arch"`
	Model    string   `json:"model"`
	Features []string `json:"features"`
	Topology `json:"topology"`
}

type Service struct {
	Host           string `json:"host"`
	ID             int    `json:"id"`
	DisabledReason string `json:"disabled_reason"`
}

type Hypervisor struct {
	// A structure that contains cpu information like arch, model, vendor, features and topology
	CpuInfo `json:"cpu_info"`
	// The current_workload is the number of tasks the hypervisor is responsible for.
	// This will be equal or greater than the number of active VMs on the system
	// (it can be greater when VMs are being deleted and the hypervisor is still cleaning up).
	CurrentWorkload int `json:"current_workload"`
	// Status of the hypervisor, either "enabled" or "disabled"
	Status string `json:"status"`
	// State of the hypervisor, either "up" or "down"
	State string `json:"state"`
	// Actual free disk on this hypervisor in GB
	DiskAvailableLeast int `json:"disk_available_least"`
	// The hypervisor's IP address
	HostIP string `json:"host_ip"`
	// The free disk remaining on this hypervisor in GB
	FreeDiskGB int `json:"free_disk_gb"`
	// The free RAM in this hypervisor in MB
	FreeRamMB int `json:"free_ram_mb"`
	// The hypervisor host name
	HypervisorHostname string `json:"hypervisor_hostname"`
	// The hypervisor type
	HypervisorType string `json:"hypervisor_type"`
	// The hypervisor version
	HypervisorVersion int `json:"hypervisor_version"`
	// Unique ID of the hypervisor
	ID int `json:"id"`
	// The disk in this hypervisor in GB
	LocalGB int `json:"local_gb"`
	// The disk used in this hypervisor in GB
	LocalGBUsed int `json:"local_gb_used"`
	// The memory of this hypervisor in MB
	MemoryMB int `json:"memory_mb"`
	// The memory used in this hypervisor in MB
	MemoryMBUsed int `json:"memory_mb_used"`
	// The number of running vms on this hypervisor
	RunningVMs int `json:"running_vms"`
	// The hypervisor service object
	Service `json:"service"`
	// The number of vcpu in this hypervisor
	VCPUs int `json:"vcpus"`
	// The number of vcpu used in this hypervisor
	VCPUsUsed int `json:"vcpus_used"`
}

func (h *Hypervisor) UnmarshalJSON(b []byte) error {

	type tmp Hypervisor
	var hypervisor *struct {
		tmp
		CpuInfo           interface{} `json:"cpu_info"`
		HypervisorVersion interface{} `json:"hypervisor_version"`
	}

	err := json.Unmarshal(b, &hypervisor)
	if err != nil {
		return err
	}

	*h = Hypervisor(hypervisor.tmp)

	// Newer versions pass the CPU into around as the correct types, this just needs
	// converting and copying into place. Older versions pass CPU info around as a string
	// and can simply be unmarshalled by the json parser
	switch t := hypervisor.CpuInfo.(type) {
	case map[string]interface{}:
		var ok bool
		if h.CpuInfo.Vendor, ok = t["vendor"].(string); !ok {
			return fmt.Errorf("CPU info vendor expected to be a string")
		}
		if h.CpuInfo.Arch, ok = t["arch"].(string); !ok {
			return fmt.Errorf("CPU info arch expected to be a string")
		}
		if h.CpuInfo.Model, ok = t["model"].(string); !ok {
			return fmt.Errorf("CPU info model expected to be a string")
		}
		if features, ok := t["features"].([]interface{}); ok {
			for _, feature := range features {
				if temp, ok := feature.(string); ok {
					h.CpuInfo.Features = append(h.CpuInfo.Features, temp)
				} else {
					return fmt.Errorf("CPU info feaure expected to be a string")
				}
			}
		} else {
			return fmt.Errorf("CPU info features to be an array")
		}
		if topology, ok := t["topology"].(map[string]interface{}); ok {
			if temp, ok := topology["sockets"].(float64); ok {
				h.CpuInfo.Topology.Sockets = int(temp)
			} else {
				return fmt.Errorf("CPU info topology sockets expected to be numeric")
			}
			if temp, ok := topology["cores"].(float64); ok {
				h.CpuInfo.Topology.Cores = int(temp)
			} else {
				return fmt.Errorf("CPU info topology cores expected to be numeric")
			}
			if temp, ok := topology["threads"].(float64); ok {
				h.CpuInfo.Topology.Threads = int(temp)
			} else {
				return fmt.Errorf("CPU info topology threads expected to be numeric")
			}
		} else {
			return fmt.Errorf("CPU info topology expected to be a map")
		}
	case string:
		err := json.Unmarshal([]byte(t), &h.CpuInfo)
		if err != nil {
			return err
		}
	}

	// A feature in OpenStack may return this value as floating point
	switch t := hypervisor.HypervisorVersion.(type) {
	case int:
		h.HypervisorVersion = t
	case float64:
		h.HypervisorVersion = int(t)
	default:
		return fmt.Errorf("Hypervisor version of unexpected type")
	}

	return nil
}

type HypervisorPage struct {
	pagination.LinkedPageBase
}

func (page HypervisorPage) IsEmpty() (bool, error) {
	hypervisors, err := ExtractHypervisors(page)
	return len(hypervisors) == 0, err
}

func ExtractHypervisors(p pagination.Page) ([]Hypervisor, error) {
	var h struct {
		Hypervisors []Hypervisor `json:"hypervisors"`
	}
	err := (p.(HypervisorPage)).ExtractInto(&h)
	return h.Hypervisors, err
}
