package extendedserverattributes

import (
	"encoding/base64"
	"encoding/json"
)

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
	Userdata *string `json:"-"`
}

func (r *ServerAttributesExt) UnmarshalJSON(b []byte) error {
	type tmp ServerAttributesExt
	var s struct {
		tmp
		Userdata *string `json:"OS-EXT-SRV-ATTR:user_data"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	*r = ServerAttributesExt(s.tmp)

	if s.Userdata != nil {
		r.Userdata = new(string)
		if v, err := base64.StdEncoding.DecodeString(*s.Userdata); err != nil {
			r.Userdata = s.Userdata
		} else {
			v := string(v)
			r.Userdata = &v
		}
	}

	return err
}
