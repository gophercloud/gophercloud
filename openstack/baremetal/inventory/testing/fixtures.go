package testing

import "github.com/gophercloud/gophercloud/openstack/baremetal/inventory"

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
	    "lldp": [],
	    "mac_address": "52:54:00:47:20:4d",
	    "name": "eth1",
	    "product": "0x0001",
	    "vendor": "0x1af4"
	},
	{
	    "client_id": null,
	    "has_carrier": true,
	    "ipv4_address": "172.24.42.100",
	    "lldp": [
		[
		    1,
		    "04112233aabbcc"
		],
		[
		    5,
		    "737730312d646973742d31622d623132"
		]
	    ],
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
			LLDP:        []inventory.LLDPTLVType{},
		},
		{
			IPV4Address: "172.24.42.100",
			MACAddress:  "52:54:00:4e:3d:30",
			Name:        "eth0",
			Product:     "0x0001",
			HasCarrier:  true,
			Vendor:      "0x1af4",
			LLDP: []inventory.LLDPTLVType{
				{
					Type:  1,
					Value: "04112233aabbcc",
				},
				{
					Type:  5,
					Value: "737730312d646973742d31622d623132",
				},
			},
			SpeedMbps: 1000,
		},
	},
	Memory: inventory.MemoryType{
		PhysicalMb: 2048.0,
		Total:      2.105864192e+09,
	},
	Hostname: "myawesomehost",
}
