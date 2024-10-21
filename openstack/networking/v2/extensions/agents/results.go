package agents

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/speakers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts an agent resource.
func (r commonResult) Extract() (*Agent, error) {
	var s struct {
		Agent *Agent `json:"agent"`
	}
	err := r.ExtractInto(&s)
	return s.Agent, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as an Agent.
type GetResult struct {
	commonResult
}

// UpdateResult represents the result of a get operation. Call its Extract
// method to interpret it as an Agent.
type UpdateResult struct {
	commonResult
}

// DeleteResult represents the result of a delete operation. Call its
// ExtractErr method to determine if the request succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// ScheduleDHCPNetworkResult represents the result of a schedule a network to
// a DHCP agent operation. ExtractErr method to determine if the request
// succeeded or failed.
type ScheduleDHCPNetworkResult struct {
	gophercloud.ErrResult
}

// RemoveDHCPNetworkResult represents the result of a remove a network from a
// DHCP agent operation. ExtractErr method to determine if the request succeeded
// or failed.
type RemoveDHCPNetworkResult struct {
	gophercloud.ErrResult
}

// ScheduleBGPSpeakerResult represents the result of adding a BGP speaker to a
// BGP DR Agent. ExtractErr method to determine if the request succeeded or
// failed.
type ScheduleBGPSpeakerResult struct {
	gophercloud.ErrResult
}

// RemoveBGPSpeakerResult represents the result of removing a BGP speaker from a
// BGP DR Agent. ExtractErr method to determine if the request succeeded or
// failed.
type RemoveBGPSpeakerResult struct {
	gophercloud.ErrResult
}

// Agent represents a Neutron agent.
type Agent struct {
	// ID is the id of the agent.
	ID string `json:"id"`

	// AdminStateUp is an administrative state of the agent.
	AdminStateUp bool `json:"admin_state_up"`

	// AgentType is a type of the agent.
	AgentType string `json:"agent_type"`

	// Alive indicates whether agent is alive or not.
	Alive bool `json:"alive"`

	// ResourcesSynced indicates whether agent is synced or not.
	// Not all agent types track resources via Placement.
	ResourcesSynced bool `json:"resources_synced"`

	// AvailabilityZone is a zone of the agent.
	AvailabilityZone string `json:"availability_zone"`

	// Binary is an executable binary of the agent.
	Binary string `json:"binary"`

	// Configurations is a configuration specific key/value pairs that are
	// determined by the agent binary and type.
	Configurations map[string]any `json:"configurations"`

	// CreatedAt is a creation timestamp.
	CreatedAt time.Time `json:"-"`

	// StartedAt is a starting timestamp.
	StartedAt time.Time `json:"-"`

	// HeartbeatTimestamp is a last heartbeat timestamp.
	HeartbeatTimestamp time.Time `json:"-"`

	// Description contains agent description.
	Description string `json:"description"`

	// Host is a hostname of the agent system.
	Host string `json:"host"`

	// Topic contains name of AMQP topic.
	Topic string `json:"topic"`
}

// UnmarshalJSON helps to convert the timestamps into the time.Time type.
func (r *Agent) UnmarshalJSON(b []byte) error {
	type tmp Agent
	var s struct {
		tmp
		CreatedAt          gophercloud.JSONRFC3339ZNoTNoZ `json:"created_at"`
		StartedAt          gophercloud.JSONRFC3339ZNoTNoZ `json:"started_at"`
		HeartbeatTimestamp gophercloud.JSONRFC3339ZNoTNoZ `json:"heartbeat_timestamp"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Agent(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.StartedAt = time.Time(s.StartedAt)
	r.HeartbeatTimestamp = time.Time(s.HeartbeatTimestamp)

	return nil
}

// AgentPage stores a single page of Agents from a List() API call.
type AgentPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of agent has
// reached the end of a page and the pager seeks to traverse over a new one.
// In order to do this, it needs to construct the next page's URL.
func (r AgentPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"agents_links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// IsEmpty determines whether or not a AgentPage is empty.
func (r AgentPage) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	agents, err := ExtractAgents(r)
	return len(agents) == 0, err
}

// ExtractAgents interprets the results of a single page from a List()
// API call, producing a slice of Agents structs.
func ExtractAgents(r pagination.Page) ([]Agent, error) {
	var s struct {
		Agents []Agent `json:"agents"`
	}
	err := (r.(AgentPage)).ExtractInto(&s)
	return s.Agents, err
}

// ListDHCPNetworksResult is the response from a List operation.
// Call its Extract method to interpret it as networks.
type ListDHCPNetworksResult struct {
	gophercloud.Result
}

// Extract interprets any ListDHCPNetworksResult as an array of networks.
func (r ListDHCPNetworksResult) Extract() ([]networks.Network, error) {
	var s struct {
		Networks []networks.Network `json:"networks"`
	}

	err := r.ExtractInto(&s)
	return s.Networks, err
}

// ListBGPSpeakersResult is the respone of agents/{id}/bgp-speakers
type ListBGPSpeakersResult struct {
	pagination.SinglePageBase
}

func (r ListBGPSpeakersResult) IsEmpty() (bool, error) {
	if r.StatusCode == 204 {
		return true, nil
	}

	speakers, err := ExtractBGPSpeakers(r)
	return 0 == len(speakers), err
}

// ExtractBGPSpeakers inteprets the ListBGPSpeakersResult into an array of BGP speakers
func ExtractBGPSpeakers(r pagination.Page) ([]speakers.BGPSpeaker, error) {
	var s struct {
		Speakers []speakers.BGPSpeaker `json:"bgp_speakers"`
	}

	err := (r.(ListBGPSpeakersResult)).ExtractInto(&s)
	return s.Speakers, err
}

// ListL3RoutersResult is the response from a List operation.
// Call its Extract method to interpret it as routers.
type ListL3RoutersResult struct {
	gophercloud.Result
}

// ScheduleL3RouterResult represents the result of a schedule a router to
// a L3 agent operation. ExtractErr method to determine if the request
// succeeded or failed.
type ScheduleL3RouterResult struct {
	gophercloud.ErrResult
}

// RemoveL3RouterResult represents the result of a remove a router from a
// L3 agent operation. ExtractErr method to determine if the request succeeded
// or failed.
type RemoveL3RouterResult struct {
	gophercloud.ErrResult
}

// Extract interprets any ListL3RoutesResult as an array of routers.
func (r ListL3RoutersResult) Extract() ([]routers.Router, error) {
	var s struct {
		Routers []routers.Router `json:"routers"`
	}

	err := r.ExtractInto(&s)
	return s.Routers, err
}
