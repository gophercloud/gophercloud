package inventory

import (
	"encoding/json"
	"fmt"
)

type ExtraDataItem map[string]interface{}

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
	var list []interface{}
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
