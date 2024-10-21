package agents

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToAgentListQuery() (string, error)
}

// ListOpts allows the filtering and sorting of paginated collections through
// the Neutron API. Filtering is achieved by passing in struct field values
// that map to the agent attributes you want to see returned.
// SortKey allows you to sort by a particular agent attribute.
// SortDir sets the direction, and is either `asc' or `desc'.
// Marker and Limit are used for the pagination.
type ListOpts struct {
	ID               string `q:"id"`
	AgentType        string `q:"agent_type"`
	Alive            *bool  `q:"alive"`
	AvailabilityZone string `q:"availability_zone"`
	Binary           string `q:"binary"`
	Description      string `q:"description"`
	Host             string `q:"host"`
	Topic            string `q:"topic"`
	Limit            int    `q:"limit"`
	Marker           string `q:"marker"`
	SortKey          string `q:"sort_key"`
	SortDir          string `q:"sort_dir"`
}

// ToAgentListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToAgentListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List returns a Pager which allows you to iterate over a collection of
// agents. It accepts a ListOpts struct, which allows you to filter and
// sort the returned collection for greater efficiency.
//
// Default policy settings return only the agents owned by the project
// of the user submitting the request, unless the user has the administrative
// role.
func List(c *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(c)
	if opts != nil {
		query, err := opts.ToAgentListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AgentPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves a specific agent based on its ID.
func Get(ctx context.Context, c *gophercloud.ServiceClient, id string) (r GetResult) {
	resp, err := c.Get(ctx, getURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToAgentUpdateMap() (map[string]any, error)
}

// UpdateOpts represents the attributes used when updating an existing agent.
type UpdateOpts struct {
	Description  *string `json:"description,omitempty"`
	AdminStateUp *bool   `json:"admin_state_up,omitempty"`
}

// ToAgentUpdateMap builds a request body from UpdateOpts.
func (opts UpdateOpts) ToAgentUpdateMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "agent")
}

// Update updates a specific agent based on its ID.
func Update(ctx context.Context, c *gophercloud.ServiceClient, id string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToAgentUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Put(ctx, updateURL(c, id), b, &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// Delete deletes a specific agent based on its ID.
func Delete(ctx context.Context, c *gophercloud.ServiceClient, id string) (r DeleteResult) {
	resp, err := c.Delete(ctx, deleteURL(c, id), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListDHCPNetworks returns a list of networks scheduled to a specific
// dhcp agent.
func ListDHCPNetworks(ctx context.Context, c *gophercloud.ServiceClient, id string) (r ListDHCPNetworksResult) {
	resp, err := c.Get(ctx, listDHCPNetworksURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ScheduleDHCPNetworkOptsBuilder allows extensions to add additional parameters
// to the ScheduleDHCPNetwork request.
type ScheduleDHCPNetworkOptsBuilder interface {
	ToAgentScheduleDHCPNetworkMap() (map[string]any, error)
}

// ScheduleDHCPNetworkOpts represents the attributes used when scheduling a
// network to a DHCP agent.
type ScheduleDHCPNetworkOpts struct {
	NetworkID string `json:"network_id" required:"true"`
}

// ToAgentScheduleDHCPNetworkMap builds a request body from ScheduleDHCPNetworkOpts.
func (opts ScheduleDHCPNetworkOpts) ToAgentScheduleDHCPNetworkMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ScheduleDHCPNetwork schedule a network to a DHCP agent.
func ScheduleDHCPNetwork(ctx context.Context, c *gophercloud.ServiceClient, id string, opts ScheduleDHCPNetworkOptsBuilder) (r ScheduleDHCPNetworkResult) {
	b, err := opts.ToAgentScheduleDHCPNetworkMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, scheduleDHCPNetworkURL(c, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveDHCPNetwork removes a network from a DHCP agent.
func RemoveDHCPNetwork(ctx context.Context, c *gophercloud.ServiceClient, id string, networkID string) (r RemoveDHCPNetworkResult) {
	resp, err := c.Delete(ctx, removeDHCPNetworkURL(c, id, networkID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListBGPSpeakers list the BGP Speakers hosted by a specific dragent
// GET /v2.0/agents/{agent-id}/bgp-drinstances
func ListBGPSpeakers(c *gophercloud.ServiceClient, agentID string) pagination.Pager {
	url := listBGPSpeakersURL(c, agentID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return ListBGPSpeakersResult{pagination.SinglePageBase(r)}
	})
}

// ScheduleBGPSpeakerOptsBuilder declare a function that build ScheduleBGPSpeakerOpts into a request body
type ScheduleBGPSpeakerOptsBuilder interface {
	ToAgentScheduleBGPSpeakerMap() (map[string]any, error)
}

// ScheduleBGPSpeakerOpts represents the data that would be POST to the endpoint
type ScheduleBGPSpeakerOpts struct {
	SpeakerID string `json:"bgp_speaker_id" required:"true"`
}

// ToAgentScheduleBGPSpeakerMap builds a request body from ScheduleBGPSpeakerOpts
func (opts ScheduleBGPSpeakerOpts) ToAgentScheduleBGPSpeakerMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ScheduleBGPSpeaker schedule a BGP speaker to a BGP agent
// POST /v2.0/agents/{agent-id}/bgp-drinstances
func ScheduleBGPSpeaker(ctx context.Context, c *gophercloud.ServiceClient, agentID string, opts ScheduleBGPSpeakerOptsBuilder) (r ScheduleBGPSpeakerResult) {
	b, err := opts.ToAgentScheduleBGPSpeakerMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, scheduleBGPSpeakersURL(c, agentID), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveBGPSpeaker removes a BGP speaker from a BGP agent
// DELETE /v2.0/agents/{agent-id}/bgp-drinstances
func RemoveBGPSpeaker(ctx context.Context, c *gophercloud.ServiceClient, agentID string, speakerID string) (r RemoveBGPSpeakerResult) {
	resp, err := c.Delete(ctx, removeBGPSpeakersURL(c, agentID, speakerID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ListDRAgentHostingBGPSpeakers the dragents that are hosting a specific bgp speaker
// GET /v2.0/bgp-speakers/{bgp-speaker-id}/bgp-dragents
func ListDRAgentHostingBGPSpeakers(c *gophercloud.ServiceClient, bgpSpeakerID string) pagination.Pager {
	url := listDRAgentHostingBGPSpeakersURL(c, bgpSpeakerID)
	return pagination.NewPager(c, url, func(r pagination.PageResult) pagination.Page {
		return AgentPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ListL3Routers returns a list of routers scheduled to a specific
// L3 agent.
func ListL3Routers(ctx context.Context, c *gophercloud.ServiceClient, id string) (r ListL3RoutersResult) {
	resp, err := c.Get(ctx, listL3RoutersURL(c, id), &r.Body, nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// ScheduleL3RouterOptsBuilder allows extensions to add additional parameters
// to the ScheduleL3Router request.
type ScheduleL3RouterOptsBuilder interface {
	ToAgentScheduleL3RouterMap() (map[string]any, error)
}

// ScheduleL3RouterOpts represents the attributes used when scheduling a
// router to a L3 agent.
type ScheduleL3RouterOpts struct {
	RouterID string `json:"router_id" required:"true"`
}

// ToAgentScheduleL3RouterMap builds a request body from ScheduleL3RouterOpts.
func (opts ScheduleL3RouterOpts) ToAgentScheduleL3RouterMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "")
}

// ScheduleL3Router schedule a router to a L3 agent.
func ScheduleL3Router(ctx context.Context, c *gophercloud.ServiceClient, id string, opts ScheduleL3RouterOptsBuilder) (r ScheduleL3RouterResult) {
	b, err := opts.ToAgentScheduleL3RouterMap()
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(ctx, scheduleL3RouterURL(c, id), b, nil, &gophercloud.RequestOpts{
		OkCodes: []int{201},
	})
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}

// RemoveL3Router removes a router from a L3 agent.
func RemoveL3Router(ctx context.Context, c *gophercloud.ServiceClient, id string, routerID string) (r RemoveL3RouterResult) {
	resp, err := c.Delete(ctx, removeL3RouterURL(c, id, routerID), nil)
	_, r.Header, r.Err = gophercloud.ParseResponse(resp, err)
	return
}
