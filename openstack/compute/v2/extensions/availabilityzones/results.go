package availabilityzones

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ServerExt is an extension to the base Server object
type ServerAvailabilityZoneExt struct {
	// AvailabilityZone is the availabilty zone the server is in.
	AvailabilityZone string `json:"OS-EXT-AZ:availability_zone"`
}

// UnmarshalJSON to override default
func (r *ServerAvailabilityZoneExt) UnmarshalJSON(b []byte) error {
	return nil
}

type StateofService struct {
	Active    bool      `json:"active"`
	Available bool      `json:"available"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UnmarshalJSON to override default
func (r *StateofService) UnmarshalJSON(b []byte) error {
	type tmp StateofService
	var s struct {
		tmp
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = StateofService(s.tmp)

	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}

type Services struct {
	NovaConductor   StateofService `json:"nova-conductor"`
	NovaConsoleauth StateofService `json:"nova-consoleauth"`
	NovaNetwork     StateofService `json:"nova-network"`
	NovaScheduler   StateofService `json:"nova-scheduler"`
	NovaCompute     StateofService `json:"nova-compute"`
}

type Hosts map[string]Services

// The current state of the availability zone
type ZoneState struct {
	// Returns true if the availability zone is available
	Available bool `json:"available"`
}

// AvailabilityZone contains all the information associated with an OpenStack
// AvailabilityZone.
type AvailabilityZone struct {
	Hosts `json:"hosts"`
	// The availability zone name
	ZoneName  string `json:"zoneName"`
	ZoneState `json:"zoneState"`
}

type OSAvailabilityZone struct {
	// The list of availability zone information
	AvailabilityZoneInfo []AvailabilityZone `json:"availabilityZoneInfo"`
}

type OSAvailabilityZonePage struct {
	pagination.SinglePageBase
}

// ExtractOSAvailabilityZones will get the OSAvailabilityZone objects out of the shareTypeAccessResult object.
func ExtractOSAvailabilityZones(r pagination.Page) (OSAvailabilityZone, error) {
	var a OSAvailabilityZone
	err := (r.(OSAvailabilityZonePage)).ExtractInto(&a)
	return a, err
}
