package clusterpolicies

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
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

// ExtractClusterPolicies provides access to the list of profiles in a page acquired from the ListDetail operation.
func ExtractClusterPolicies(r pagination.Page) ([]ClusterPolicy, error) {
	var s struct {
		ClusterPolicies []ClusterPolicy `json:"cluster_policies"`
	}
	err := (r.(ClusterPolicyPage)).ExtractInto(&s)
	return s.ClusterPolicies, err
}

// ClusterPolicyPage contains a single page of all policies from a ListDetails call.
type ClusterPolicyPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines if a ProfilePage contains any results.
func (page ClusterPolicyPage) IsEmpty() (bool, error) {
	clusterPolicies, err := ExtractClusterPolicies(page)
	return len(clusterPolicies) == 0, err
}
