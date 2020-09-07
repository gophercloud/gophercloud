package quotas

import "github.com/gophercloud/gophercloud"

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a Quota resource.
func (r commonResult) Extract() (*Quota, error) {
	var s struct {
		Quota *Quota `json:"quota"`
	}
	err := r.ExtractInto(&s)
	return s.Quota, err
}

// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Quota.
type GetResult struct {
	commonResult
}

// Quota contains load balancer quotas for a project.
type Quota struct {
	// Loadbalancer represents the number of load balancers. A "-1" value means no limit.
	Loadbalancer int `json:"loadbalancer"`

	// Listener represents the number of listeners. A "-1" value means no limit.
	Listener int `json:"listener"`

	// Member represents the number of members. A "-1" value means no limit.
	Member int `json:"member"`

	// Poool represents the number of pools. A "-1" value means no limit.
	Pool int `json:"pool"`

	// HealthMonitor represents the number of healthmonitors. A "-1" value means no limit.
	Healthmonitor int `json:"healthmonitor"`

	// L7Policy represents the number of l7policies. A "-1" value means no limit.
	L7Policy int `json:"l7policy"`

	// L7Rule represents the number of l7rules. A "-1" value means no limit.
	L7Rule int `json:"l7rule"`
}
