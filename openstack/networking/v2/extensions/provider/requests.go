package provider

import (
	"github.com/gophercloud/gophercloud"
)

type CreateOptsExtBuilder interface {
	ToProviderCreateMap() (map[string]interface{}, error)
}

type CreateOptsExt struct {
	AdminStateUp    *bool     `json:"admin_state_up,omitempty"`
	Name            string    `json:"name,omitempty"`
	Shared          *bool     `json:"shared,omitempty"`
	TenantID        string    `json:"tenant_id,omitempty"`
	NetworkType     string    `json:"provider:network_type,omitempty"`
	PhysicalNetwork string    `json:"provider:physical_network,omitempty"`
	SegmentationID  string    `json:"provider:segmentation_id,omitempty"`
	Segments        []Segment `json:"segments,omitempty"`
}

func (opts CreateOptsExt) ToProviderCreateMap() (map[string]interface{}, error) {
	return gophercloud.BuildRequestBody(opts, "network")
}

func Create(c *gophercloud.ServiceClient, opts CreateOptsExtBuilder) (r CreateResult) {
	b, err := opts.ToProviderCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Post(createURL(c), b, &r.Body, nil)
	return
}
