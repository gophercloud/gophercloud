package speakers

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

const jroot = "bgp_speaker"

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a bgp speaker resource.
func (r commonResult) Extract() (*BGPSpeaker, error) {
	var s BGPSpeaker
	err := r.ExtractInto(&s)
	return &s, err
}

func (r commonResult) ExtractInto(v interface{}) error {
	return r.Result.ExtractIntoStructPtr(v, jroot)
}

// BGPSpeaker BGP Speaker
type BGPSpeaker struct {
	// UUID for the bgp speaker
	ID string `json:"id"`

	// Human-readable name for the bgp speaker. Might not be unique.
	Name string `json:"name"`

	// TenantID is the project owner of the bgp speaker.
	TenantID string `json:"tenant_id"`

	// ProjectID is the project owner of the bgp speaker.
	ProjectID string `json:"project_id"`

	// If the speaker would advertise floating ip host routes
	AdvertiseFloatingIPHostRoutes bool `json:"advertise_floating_ip_host_routes"`

	// If the speaker would advertise tenant networks
	AdvertiseTenantNetworks bool `json:"advertise_tenant_networks"`

	// IP version
	IPVersion int `json:"ip_version"`

	// Local Autonomous System
	LocalAS int `json:"local_as"`

	// The uuid of the Networks configured with this speaker
	Networks []string `json:"networks"`

	// The uuid of the BGP Peer Configured with this speaker
	Peers []string `json:"peers"`
}

// BGPSpeakerPage is the page returned by a pager when traversing over a
// collection of bgp speakers.
type BGPSpeakerPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether a BGPSpeakerPage struct is empty.
func (r BGPSpeakerPage) IsEmpty() (bool, error) {
	is, err := ExtractBGPSpeakers(r)
	return len(is) == 0, err
}

// ExtractBGPSpeakers accepts a Page struct, specifically a BGPSpeakerPage struct,
// and extracts the elements into a slice of BGPSpeaker structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractBGPSpeakers(r pagination.Page) ([]BGPSpeaker, error) {
	var s []BGPSpeaker
	err := ExtractBGPSpeakersInto(r, &s)
	return s, err
}

func ExtractBGPSpeakersInto(r pagination.Page, v interface{}) error {
	return r.(BGPSpeakerPage).Result.ExtractIntoSlicePtr(v, "bgp_speakers")
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a BGPSpeaker.
type GetResult struct {
	commonResult
}
