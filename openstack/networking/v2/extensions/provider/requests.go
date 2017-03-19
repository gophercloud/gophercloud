package provider

import (
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

type CreateOptsExtBuilder interface {
	ToProviderCreateMap() (map[string]interface{}, error)
}

type CreateOptsExt struct {
	networks.CreateOptsBuilder
	Segments        []Segment `json:"segments,omitempty"`
}

func (opts CreateOptsExt) ToNetworkCreateMap() (map[string]interface{}, error) {
	base, err := opts.CreateOptsBuilder.ToNetworkCreateMap()
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
