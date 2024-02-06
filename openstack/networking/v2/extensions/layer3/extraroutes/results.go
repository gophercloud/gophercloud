package extraroutes

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
)

// Extract is a function that accepts a result and extracts a router.
func (r commonResult) Extract() (*routers.Router, error) {
	var s struct {
		Router *routers.Router `json:"router"`
	}
	err := r.ExtractInto(&s)
	return s.Router, err
}

type commonResult struct {
	gophercloud.Result
}

// AddResult represents the result of an extra routes add operation. Call its
// Extract method to interpret it as a *routers.Router.
type AddResult struct {
	commonResult
}

// RemoveResult represents the result of an extra routes remove operation. Call
// its Extract method to interpret it as a *routers.Router.
type RemoveResult struct {
	commonResult
}
