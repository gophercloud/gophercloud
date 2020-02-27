package extendedserverattributes

// ServerAttributesExt represents basic OS-EXT-SRV-ATTR server response fields.
// You should use extract methods from microversions.go to retrieve additional
// fields.
type ServerAttributesExt struct {
	// Host is the host/hypervisor that the instance is hosted on.
	Host string `json:"OS-EXT-SRV-ATTR:host"`

	// InstanceName is the name of the instance.
	InstanceName string `json:"OS-EXT-SRV-ATTR:instance_name"`

	// HypervisorHostname is the hostname of the host/hypervisor that the
	// instance is hosted on.
	HypervisorHostname string `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`

	// ReservationID is the reservation ID of the instance.
	// This requires microversion 2.3 or later.
	ReservationID *string `json:"OS-EXT-SRV-ATTR:reservation_id"`

	// LaunchIndex is the launch index of the instance.
	// This requires microversion 2.3 or later.
	LaunchIndex *int `json:"OS-EXT-SRV-ATTR:launch_index"`

	// RAMDiskID is the ID of the RAM disk image of the instance.
	// This requires microversion 2.3 or later.
	RAMDiskID *string `json:"OS-EXT-SRV-ATTR:ramdisk_id"`

	// KernelID is the ID of the kernel image of the instance.
	// This requires microversion 2.3 or later.
	KernelID *string `json:"OS-EXT-SRV-ATTR:kernel_id"`

	// Hostname is the hostname of the instance.
	// This requires microversion 2.3 or later.
	Hostname *string `json:"OS-EXT-SRV-ATTR:hostname"`

	// RootDeviceName is the name of the root device of the instance.
	// This requires microversion 2.3 or later.
	RootDeviceName *string `json:"OS-EXT-SRV-ATTR:root_device_name"`

	// Userdata is the userdata of the instance.
	// This requires microversion 2.3 or later.
	Userdata *string `json:"OS-EXT-SRV-ATTR:userdata"`
}
