package testing

import (
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/baremetal/inventory"
)

const InventorySample = `{
    "bmc_address": "192.167.2.134",
    "boot": {
	"current_boot_mode": "bios",
	"pxe_interface": "52:54:00:4e:3d:30"
    },
    "cpu": {
	"architecture": "x86_64",
	"count": 2,
	"flags": [
	    "fpu",
	    "mmx",
	    "fxsr",
	    "sse",
	    "sse2"
	],
	"frequency": "2100.084"
    },
    "disks": [
	{
	    "hctl": null,
	    "model": "",
	    "name": "/dev/vda",
	    "rotational": true,
	    "serial": null,
	    "size": 13958643712,
	    "vendor": "0x1af4",
	    "wwn": null,
	    "wwn_vendor_extension": null,
	    "wwn_with_extension": null
	}
    ],
    "hostname": "myawesomehost",
    "interfaces": [
	{
	    "client_id": null,
	    "has_carrier": true,
	    "ipv4_address": "172.24.42.101",
	    "mac_address": "52:54:00:47:20:4d",
	    "name": "eth1",
	    "product": "0x0001",
	    "vendor": "0x1af4"
	},
	{
	    "client_id": null,
	    "has_carrier": true,
	    "ipv4_address": "172.24.42.100",
	    "mac_address": "52:54:00:4e:3d:30",
	    "name": "eth0",
	    "product": "0x0001",
	    "vendor": "0x1af4",
	    "speed_mbps": 1000
	}
    ],
    "memory": {
	"physical_mb": 2048,
	"total": 2105864192
    },
    "system_vendor": {
	"manufacturer": "Bochs",
	"product_name": "Bochs",
	"serial_number": "Not Specified",
	"firmware": {
	    "version": "1.2.3.4"
	}
    }
}`

// ExtraDataJSONSample contains extra hardware sample data reported by the inspection process.
const ExtraDataJSONSample = `
{
  "cpu": {
    "logical": {
      "number": 16
    },
    "physical": {
      "clock": 2105032704,
      "cores": 8,
      "flags": "lm fpu fpu_exception wp vme de"
    }
  },
  "disk": {
    "sda": {
      "rotational": 1,
      "vendor": "TEST"
    }
  },
  "firmware": {
    "bios": {
      "date": "01/01/1970",
      "vendor": "test"
    }
  },
  "ipmi": {
    "Fan1A RPM": {
      "unit": "RPM",
      "value": 3120
    },
    "Fan1B RPM": {
      "unit": "RPM",
      "value": 2280
    }
  },
  "memory": {
    "bank0": {
      "clock": 1600000000.0,
      "description": "DIMM DDR3 Synchronous Registered (Buffered) 1600 MHz (0.6 ns)"
    },
    "bank1": {
      "clock": 1600000000.0,
      "description": "DIMM DDR3 Synchronous Registered (Buffered) 1600 MHz (0.6 ns)"
    }
  },
  "network": {
    "em1": {
      "Autonegotiate": "on",
      "loopback": "off [fixed]"
    },
    "p2p1": {
      "Autonegotiate": "on",
      "loopback": "off [fixed]"
    }
  },
  "system": {
    "ipmi": {
      "channel": 1
    },
    "kernel": {
      "arch": "x86_64",
      "version": "3.10.0"
    },
    "motherboard": {
      "vendor": "Test"
    },
    "product": {
      "name": "test",
      "vendor": "Test"
    }
  }
}
`

// NUMADataJSONSample contains NUMA sample data reported by the inspection process.
const NUMADataJSONSample = `
{
  "numa_topology": {
    "cpus": [
      {
        "cpu": 6,
        "numa_node": 1,
        "thread_siblings": [
          3,
          27
        ]
      },
      {
        "cpu": 10,
        "numa_node": 0,
        "thread_siblings": [
          20,
          44
        ]
      }
    ],
    "nics": [
      {
        "name": "p2p1",
        "numa_node": 0
      },
      {
        "name": "p2p2",
        "numa_node": 1
      }
    ],
    "ram": [
      {
        "numa_node": 0,
        "size_kb": 99289532
      },
      {
        "numa_node": 1,
        "size_kb": 100663296
      }
    ]
  }
}
`

var StandardPluginDataSample = fmt.Sprintf(`
{
    "all_interfaces": {
        "eth0": {
	    "client_id": null,
	    "has_carrier": true,
	    "ipv4_address": "172.24.42.101",
	    "mac_address": "52:54:00:47:20:4d",
	    "name": "eth1",
	    "product": "0x0001",
	    "vendor": "0x1af4",
	    "pxe_enabled": true
	},
        "eth1": {
	    "client_id": null,
	    "has_carrier": true,
	    "ipv4_address": "172.24.42.100",
	    "mac_address": "52:54:00:4e:3d:30",
	    "name": "eth0",
	    "product": "0x0001",
	    "vendor": "0x1af4",
	    "speed_mbps": 1000,
	    "pxe_enabled": false
	}
    },
    "boot_interface": "52:54:00:4e:3d:30",
    "configuration": {
	"collectors": ["default", "logs"],
	"managers": [
	    {
		"name": "generic_hardware_manager",
		"version": "1.1"
	    }
	]
    },
    "error": null,
    "extra": %s,
    "valid_interfaces": {
        "eth0": {
	    "client_id": null,
	    "has_carrier": true,
	    "ipv4_address": "172.24.42.101",
	    "mac_address": "52:54:00:47:20:4d",
	    "name": "eth1",
	    "product": "0x0001",
	    "vendor": "0x1af4",
	    "pxe_enabled": true
	}
    },
    "lldp_raw": {
	"eth0": [
	    [
		1,
		"04112233aabbcc"
	    ],
	    [
		5,
		"737730312d646973742d31622d623132"
	    ]
	]
    },
    "macs": [
        "52:54:00:4e:3d:30"
    ],
    "parsed_lldp": {
        "eth0": {
            "switch_chassis_id": "11:22:33:aa:bb:cc",
            "switch_system_name": "sw01-dist-1b-b12"
        }
    },
    "root_disk": {
        "hctl": null,
        "model": "",
        "name": "/dev/vda",
        "rotational": true,
        "serial": null,
        "size": 13958643712,
        "vendor": "0x1af4",
        "wwn": null,
        "wwn_vendor_extension": null,
        "wwn_with_extension": null
    }
}`, ExtraDataJSONSample)

var Inventory = inventory.InventoryType{
	SystemVendor: inventory.SystemVendorType{
		Manufacturer: "Bochs",
		ProductName:  "Bochs",
		SerialNumber: "Not Specified",
		Firmware: inventory.SystemFirmwareType{
			Version: "1.2.3.4",
		},
	},
	BmcAddress: "192.167.2.134",
	Boot: inventory.BootInfoType{
		CurrentBootMode: "bios",
		PXEInterface:    "52:54:00:4e:3d:30",
	},
	CPU: inventory.CPUType{
		Count:        2,
		Flags:        []string{"fpu", "mmx", "fxsr", "sse", "sse2"},
		Frequency:    "2100.084",
		Architecture: "x86_64",
	},
	Disks: []inventory.RootDiskType{
		{
			Rotational: true,
			Model:      "",
			Name:       "/dev/vda",
			Size:       13958643712,
			Vendor:     "0x1af4",
		},
	},
	Interfaces: []inventory.InterfaceType{
		{
			Vendor:      "0x1af4",
			HasCarrier:  true,
			MACAddress:  "52:54:00:47:20:4d",
			Name:        "eth1",
			Product:     "0x0001",
			IPV4Address: "172.24.42.101",
		},
		{
			IPV4Address: "172.24.42.100",
			MACAddress:  "52:54:00:4e:3d:30",
			Name:        "eth0",
			Product:     "0x0001",
			HasCarrier:  true,
			Vendor:      "0x1af4",
			SpeedMbps:   1000,
		},
	},
	Memory: inventory.MemoryType{
		PhysicalMb: 2048.0,
		Total:      2.105864192e+09,
	},
	Hostname: "myawesomehost",
}

var ExtraData = inventory.ExtraDataType{
	CPU: inventory.ExtraDataSection{
		"logical": map[string]any{
			"number": float64(16),
		},
		"physical": map[string]any{
			"clock": float64(2105032704),
			"cores": float64(8),
			"flags": "lm fpu fpu_exception wp vme de",
		},
	},
	Disk: inventory.ExtraDataSection{
		"sda": map[string]any{
			"rotational": float64(1),
			"vendor":     "TEST",
		},
	},
	Firmware: inventory.ExtraDataSection{
		"bios": map[string]any{
			"date":   "01/01/1970",
			"vendor": "test",
		},
	},
	IPMI: inventory.ExtraDataSection{
		"Fan1A RPM": map[string]any{
			"unit":  "RPM",
			"value": float64(3120),
		},
		"Fan1B RPM": map[string]any{
			"unit":  "RPM",
			"value": float64(2280),
		},
	},
	Memory: inventory.ExtraDataSection{
		"bank0": map[string]any{
			"clock":       1600000000.0,
			"description": "DIMM DDR3 Synchronous Registered (Buffered) 1600 MHz (0.6 ns)",
		},
		"bank1": map[string]any{
			"clock":       1600000000.0,
			"description": "DIMM DDR3 Synchronous Registered (Buffered) 1600 MHz (0.6 ns)",
		},
	},
	Network: inventory.ExtraDataSection{
		"em1": map[string]any{
			"Autonegotiate": "on",
			"loopback":      "off [fixed]",
		},
		"p2p1": map[string]any{
			"Autonegotiate": "on",
			"loopback":      "off [fixed]",
		},
	},
	System: inventory.ExtraDataSection{
		"ipmi": map[string]any{
			"channel": float64(1),
		},
		"kernel": map[string]any{
			"arch":    "x86_64",
			"version": "3.10.0",
		},
		"motherboard": map[string]any{
			"vendor": "Test",
		},
		"product": map[string]any{
			"name":   "test",
			"vendor": "Test",
		},
	},
}

var NUMATopology = inventory.NUMATopology{
	CPUs: []inventory.NUMACPU{
		{
			CPU:            6,
			NUMANode:       1,
			ThreadSiblings: []int{3, 27},
		},
		{
			CPU:            10,
			NUMANode:       0,
			ThreadSiblings: []int{20, 44},
		},
	},
	NICs: []inventory.NUMANIC{
		{
			Name:     "p2p1",
			NUMANode: 0,
		},
		{
			Name:     "p2p2",
			NUMANode: 1,
		},
	},
	RAM: []inventory.NUMARAM{
		{
			NUMANode: 0,
			SizeKB:   99289532,
		},
		{
			NUMANode: 1,
			SizeKB:   100663296,
		},
	},
}

var StandardPluginData = inventory.StandardPluginData{
	AllInterfaces: map[string]inventory.ProcessedInterfaceType{
		"eth0": {
			InterfaceType: inventory.InterfaceType{
				Vendor:      "0x1af4",
				HasCarrier:  true,
				MACAddress:  "52:54:00:47:20:4d",
				Name:        "eth1",
				Product:     "0x0001",
				IPV4Address: "172.24.42.101",
			},
			PXEEnabled: true,
		},
		"eth1": {
			InterfaceType: inventory.InterfaceType{
				IPV4Address: "172.24.42.100",
				MACAddress:  "52:54:00:4e:3d:30",
				Name:        "eth0",
				Product:     "0x0001",
				HasCarrier:  true,
				Vendor:      "0x1af4",
				SpeedMbps:   1000,
			},
		},
	},
	BootInterface: "52:54:00:4e:3d:30",
	Configuration: inventory.ConfigurationType{
		Collectors: []string{"default", "logs"},
		Managers: []inventory.HardwareManager{
			{
				Name:    "generic_hardware_manager",
				Version: "1.1",
			},
		},
	},
	Error: "",
	Extra: ExtraData,
	MACs:  []string{"52:54:00:4e:3d:30"},
	ParsedLLDP: map[string]inventory.ParsedLLDP{
		"eth0": map[string]any{
			"switch_chassis_id":  "11:22:33:aa:bb:cc",
			"switch_system_name": "sw01-dist-1b-b12",
		},
	},
	RawLLDP: map[string][]inventory.LLDPTLVType{
		"eth0": {
			{
				Type:  1,
				Value: "04112233aabbcc",
			},
			{
				Type:  5,
				Value: "737730312d646973742d31622d623132",
			},
		},
	},
	RootDisk: inventory.RootDiskType{
		Rotational: true,
		Model:      "",
		Name:       "/dev/vda",
		Size:       13958643712,
		Vendor:     "0x1af4",
	},
	ValidInterfaces: map[string]inventory.ProcessedInterfaceType{
		"eth0": {
			InterfaceType: inventory.InterfaceType{
				Vendor:      "0x1af4",
				HasCarrier:  true,
				MACAddress:  "52:54:00:47:20:4d",
				Name:        "eth1",
				Product:     "0x0001",
				IPV4Address: "172.24.42.101",
			},
			PXEEnabled: true,
		},
	},
}
