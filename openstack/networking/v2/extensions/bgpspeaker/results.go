package bgpspeaker

import (
	"encoding/json"
)

// BGP Speaker
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


func (n *BGPSpeaker) UnmarshalJSON(b []byte) error {
        type tmp BGPSpeaker
        var bgpspeaker struct {
            tmp
        }
        if err := json.Unmarshal(b, &bgpspeaker); err != nil {
            return err
        }
        *n = BGPSpeaker(bgpspeaker.tmp)
        return nil
}
