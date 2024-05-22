package inventory

import (
	"encoding/json"
	"fmt"
)

type ExtraDataItem map[string]any

type ExtraDataSection map[string]ExtraDataItem

type ExtraDataType struct {
	CPU      ExtraDataSection `json:"cpu"`
	Disk     ExtraDataSection `json:"disk"`
	Firmware ExtraDataSection `json:"firmware"`
	IPMI     ExtraDataSection `json:"ipmi"`
	Memory   ExtraDataSection `json:"memory"`
	Network  ExtraDataSection `json:"network"`
	System   ExtraDataSection `json:"system"`
}

type NUMATopology struct {
	CPUs []NUMACPU `json:"cpus"`
	NICs []NUMANIC `json:"nics"`
	RAM  []NUMARAM `json:"ram"`
}

type NUMACPU struct {
	CPU            int   `json:"cpu"`
	NUMANode       int   `json:"numa_node"`
	ThreadSiblings []int `json:"thread_siblings"`
}

type NUMANIC struct {
	Name     string `json:"name"`
	NUMANode int    `json:"numa_node"`
}

type NUMARAM struct {
	NUMANode int `json:"numa_node"`
	SizeKB   int `json:"size_kb"`
}

type LLDPTLVType struct {
	Type  int
	Value string
}

// UnmarshalJSON interprets an LLDP TLV [key, value] pair as an LLDPTLVType structure
func (r *LLDPTLVType) UnmarshalJSON(data []byte) error {
	var list []any
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}

	if len(list) != 2 {
		return fmt.Errorf("Invalid LLDP TLV key-value pair")
	}

	fieldtype, ok := list[0].(float64)
	if !ok {
		return fmt.Errorf("LLDP TLV key is not number")
	}

	value, ok := list[1].(string)
	if !ok {
		return fmt.Errorf("LLDP TLV value is not string")
	}

	r.Type = int(fieldtype)
	r.Value = value
	return nil
}

type HardwareManager struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ConfigurationType struct {
	// Collectors is a list of enabled collectors - ramdisk-side inspection
	// plugins that populated the plugin data.
	Collectors []string `json:"collectors"`
	// Managers is a list of hardware managers - ramdisk-side plugins that
	// implement all actions, such as writing images or collecting
	// inventory.
	Managers []HardwareManager `json:"managers"`
}

type ParsedLLDP = map[string]any

type ProcessedInterfaceType struct {
	InterfaceType
	// Whether PXE was enabled on this interface during inspection
	PXEEnabled bool `json:"pxe_enabled"`
}

// StandardPluginData represents the plugin data as collected and processes
// by a standard ramdisk and a standard Ironic deployment.
// The format and contents of the stored data depends on the ramdisk used
// and plugins enabled both in the ramdisk and in inspector itself.
// This structure has been provided for basic compatibility but it
// will need extensions.
type StandardPluginData struct {
	AllInterfaces   map[string]ProcessedInterfaceType `json:"all_interfaces"`
	BootInterface   string                            `json:"boot_interface"`
	Configuration   ConfigurationType                 `json:"configuration"`
	Error           string                            `json:"error"`
	Extra           ExtraDataType                     `json:"extra"`
	MACs            []string                          `json:"macs"`
	NUMATopology    NUMATopology                      `json:"numa_topology"`
	ParsedLLDP      map[string]ParsedLLDP             `json:"parsed_lldp"`
	RawLLDP         map[string][]LLDPTLVType          `json:"lldp_raw"`
	RootDisk        RootDiskType                      `json:"root_disk"`
	ValidInterfaces map[string]ProcessedInterfaceType `json:"valid_interfaces"`
}
