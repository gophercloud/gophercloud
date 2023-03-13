package availabilityzones

import (
	"github.com/gophercloud/gophercloud/pagination"
)

// AvailabilityZoneState contains the state information associated with an
// OpenStack Block Storage AvailabilityZone.
type AvailabilityZoneState struct {
	// Whether the availability zone is available.
	Available bool `json:"available"`
}

// AvailabilityZone contains all the information associated with an OpenStack
// Block Storage AvailabilityZone.
type AvailabilityZone struct {
	// The name of the availability zone.
	Name string `json:"zoneName"`
	// The state of the availability zone.
	State AvailabilityZoneState `json:"zoneState"`
}

// AvailabilityZonePage contains a single page of all Availability Zones.
type AvailabilityZonePage struct {
	pagination.SinglePageBase
}

// ExtractAvailabilityZones will remove the envelope from the response and get
// the AvailabilityZone objects.
func ExtractAvailabilityZones(r pagination.Page) ([]AvailabilityZone, error) {
	var a struct {
		AvailabilityZones []AvailabilityZone `json:"availabilityZoneInfo"`
	}
	err := (r.(AvailabilityZonePage)).ExtractInto(&a)
	return a.AvailabilityZones, err
}
