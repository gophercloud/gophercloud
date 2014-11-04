package vips

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

// VIP represents a Virtual IP API resource.
type VIP struct {
	Address string `json:"address,omitempty"`
	ID      int    `json:"id,omitempty"`
	Type    string `json:"type,omitempty"`
	Version string `json:"ipVersion,omitempty" mapstructure:"ipVersion"`
}

// VIPPage is the page returned by a pager when traversing over a collection
// of VIPs.
type VIPPage struct {
	pagination.SinglePageBase
}

// IsEmpty checks whether a VIPPage struct is empty.
func (p VIPPage) IsEmpty() (bool, error) {
	is, err := ExtractVIPs(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractVIPs accepts a Page struct, specifically a VIPPage struct, and
// extracts the elements into a slice of VIP structs. In other words, a
// generic collection is mapped into a relevant slice.
func ExtractVIPs(page pagination.Page) ([]VIP, error) {
	var resp struct {
		VIPs []VIP `mapstructure:"virtualIps" json:"virtualIps"`
	}

	err := mapstructure.Decode(page.(VIPPage).Body, &resp)

	return resp.VIPs, err
}
