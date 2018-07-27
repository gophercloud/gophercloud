package extendedserverattributes

// ServerAttributesExt represents OS-EXT-SRV-ATTR server response fields.
type ServerAttributesExt struct {
	ReservationID      string `json:"OS-EXT-SRV-ATTR:reservation_id"`
	LaunchIndex        int    `json:"OS-EXT-SRV-ATTR:launch_index"`
	Hostname           string `json:"OS-EXT-SRV-ATTR:hostname"`
	Host               string `json:"OS-EXT-SRV-ATTR:host"`
	KernelID           string `json:"OS-EXT-SRV-ATTR:kernel_id"`
	RamdiskID          string `json:"OS-EXT-SRV-ATTR:ramdisk_id"`
	RootDeviceName     string `json:"OS-EXT-SRV-ATTR:root_device_name"`
	UserData           string `json:"OS-EXT-SRV-ATTR:user_data"`
	InstanceName       string `json:"OS-EXT-SRV-ATTR:instance_name"`
	HypervisorHostname string `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
}
