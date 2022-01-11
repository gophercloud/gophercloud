package provider

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

// CreateOptsExt adds a Segments option to the base Network CreateOpts.
type CreateOptsExt struct {
	networks.CreateOptsBuilder
	NetworkType     string    `json:"provider:network_type,omitempty"`
	PhysicalNetwork string    `json:"provider:physical_network,omitempty"`
	SegmentationID  int       `json:"provider:segmentation_id,omitempty"`
	Segments        []Segment `json:"segments,omitempty"`
}

// ToNetworkCreateMap adds segments to the base network creation options.
func (opts CreateOptsExt) ToNetworkCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOptsBuilder.ToNetworkCreateMap()
	if err != nil {
		return nil, err
	}

	network := base["network"].(map[string]interface{})

	if opts.NetworkType != "" {
		network["provider:network_type"] = opts.NetworkType
	}

	if opts.PhysicalNetwork != "" {
		network["provider:physical_network"] = opts.PhysicalNetwork
	}

	if opts.SegmentationID > 0 {
		network["provider:segmentation_id"] = opts.SegmentationID
	}

	if opts.Segments != nil {
		network["segments"] = opts.Segments
	}

	return base, nil
}

// UpdateOptsExt adds a Segments option to the base Network UpdateOpts.
type UpdateOptsExt struct {
	networks.UpdateOptsBuilder
	Segments *[]Segment `json:"segments,omitempty"`
}

// ToNetworkUpdateMap adds segments to the base network update options.
func (opts UpdateOptsExt) ToNetworkUpdateMap() (map[string]interface{}, error) {
	base, err := opts.UpdateOptsBuilder.ToNetworkUpdateMap()
	if err != nil {
		return nil, err
	}

	if opts.Segments == nil {
		return base, nil
	}

	providerMap := base["network"].(map[string]interface{})
	providerMap["segments"] = opts.Segments

	return base, nil
}
