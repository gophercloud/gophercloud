package clusterpolicies

import (
	"github.com/gophercloud/gophercloud"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// Extract provides access to the individual Policy returned by the Get and
// Create functions.
func (r commonResult) Extract() (*ClusterPolicy, error) {
	var s struct {
		ClusterPolicy *ClusterPolicy `json:"cluster_policy"`
	}
	err := r.ExtractInto(&s)
	return s.ClusterPolicy, err
}

type ClusterPolicy struct {
	ClusterUUID string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Enabled     bool   `json:"enabled"`
	ID          string `json:"id"`
	PolicyID    string `json:"policy_id"`
	PolicyName  string `json:"policy_name"`
	PolicyType  string `json:"policy_type"`
}
