package inventory

type BootInfoType struct {
	CurrentBootMode string `json:"current_boot_mode"`
	PXEInterface    string `json:"pxe_interface"`
}

type CPUType struct {
	Architecture string   `json:"architecture"`
	Count        int      `json:"count"`
	Flags        []string `json:"flags"`
	Frequency    string   `json:"frequency"`
	ModelName    string   `json:"model_name"`
}

type InterfaceType struct {
	BIOSDevName string `json:"biosdevname"`
	ClientID    string `json:"client_id"`
	HasCarrier  bool   `json:"has_carrier"`
	IPV4Address string `json:"ipv4_address"`
	IPV6Address string `json:"ipv6_address"`
	MACAddress  string `json:"mac_address"`
	Name        string `json:"name"`
	Product     string `json:"product"`
	SpeedMbps   int    `json:"speed_mbps"`
	Vendor      string `json:"vendor"`
}

type MemoryType struct {
	PhysicalMb int `json:"physical_mb"`
	Total      int `json:"total"`
}

type RootDiskType struct {
	Hctl               string `json:"hctl"`
	Model              string `json:"model"`
	Name               string `json:"name"`
	ByPath             string `json:"by_path"`
	Rotational         bool   `json:"rotational"`
	Serial             string `json:"serial"`
	Size               int64  `json:"size"`
	Vendor             string `json:"vendor"`
	Wwn                string `json:"wwn"`
	WwnVendorExtension string `json:"wwn_vendor_extension"`
	WwnWithExtension   string `json:"wwn_with_extension"`
}

type SystemFirmwareType struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	Vendor    string `json:"vendor"`
}

type SystemVendorType struct {
	Manufacturer string             `json:"manufacturer"`
	ProductName  string             `json:"product_name"`
	SerialNumber string             `json:"serial_number"`
	Firmware     SystemFirmwareType `json:"firmware"`
}

type InventoryType struct {
	BmcAddress   string           `json:"bmc_address"`
	Boot         BootInfoType     `json:"boot"`
	CPU          CPUType          `json:"cpu"`
	Disks        []RootDiskType   `json:"disks"`
	Interfaces   []InterfaceType  `json:"interfaces"`
	Memory       MemoryType       `json:"memory"`
	SystemVendor SystemVendorType `json:"system_vendor"`
	Hostname     string           `json:"hostname"`
}
