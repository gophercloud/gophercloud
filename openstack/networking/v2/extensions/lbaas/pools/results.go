package pools

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud/pagination"
)

type Pool struct {
	// The status of the pool. Indicates whether the pool is operational.
	Status string

	// The load-balancer algorithm, which is round-robin, least-connections, and
	// so on. This value, which must be supported, is dependent on the provider.
	// Round-robin must be supported.
	LBMethod string `json:"lb_method" mapstructure:"lb_method"`

	// The protocol of the pool, which is TCP, HTTP, or HTTPS.
	Protocol string

	// Description for the pool.
	Description string

	// The IDs of associated monitors which check the health of the pool members.
	MonitorIDs []string `json:"health_monitors" mapstructure:"health_monitors"`

	// The ID of the subnet associated with the pool
	SubnetID string `json:"subnet_id" mapstructure:"subnet_id"`

	// Owner of the pool. Only an administrative user can specify a tenant ID
	// other than its own.
	TenantID string `json:"tenant_id" mapstructure:"tenant_id"`

	// The administrative state of the pool, which is up (true) or down (false).
	AdminStateUp bool `json:"admin_state_up" mapstructure:"admin_state_up"`

	// Pool name. Does not have to be unique.
	Name string

	// List of member IDs that belong to the pool.
	MemberIDs []string `json:"members" mapstructure:"members"`

	// The unique ID for the pool.
	ID string

	// The ID of the virtual IP associated with this pool
	VIPID string `json:"vip_id" mapstructure:"vip_id"`
}

// PoolPage is the page returned by a pager when traversing over a
// collection of pools.
type PoolPage struct {
	pagination.LinkedPageBase
}

// NextPageURL is invoked when a paginated collection of pools has reached
// the end of a page and the pager seeks to traverse over a new one. In order
// to do this, it needs to construct the next page's URL.
func (p PoolPage) NextPageURL() (string, error) {
	type link struct {
		Href string `mapstructure:"href"`
		Rel  string `mapstructure:"rel"`
	}
	type resp struct {
		Links []link `mapstructure:"pools_links"`
	}

	var r resp
	err := mapstructure.Decode(p.Body, &r)
	if err != nil {
		return "", err
	}

	var url string
	for _, l := range r.Links {
		if l.Rel == "next" {
			url = l.Href
		}
	}
	if url == "" {
		return "", nil
	}

	return url, nil
}

// IsEmpty checks whether a PoolPage struct is empty.
func (p PoolPage) IsEmpty() (bool, error) {
	is, err := ExtractPools(p)
	if err != nil {
		return true, nil
	}
	return len(is) == 0, nil
}

// ExtractRouters accepts a Page struct, specifically a RouterPage struct,
// and extracts the elements into a slice of Router structs. In other words,
// a generic collection is mapped into a relevant slice.
func ExtractPools(page pagination.Page) ([]Pool, error) {
	var resp struct {
		Pools []Pool `mapstructure:"pools" json:"pools"`
	}

	err := mapstructure.Decode(page.(PoolPage).Body, &resp)
	if err != nil {
		return nil, err
	}

	return resp.Pools, nil
}
