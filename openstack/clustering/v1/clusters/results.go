package clusters

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response of a Create operations.
type CreateResult struct {
	commonResult
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// UpdateResult is the response of a Update operations.
type UpdateResult struct {
	commonResult
}

// DeleteResult is the result from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// ResizeResult is the response of a Get operations.
type ResizeResult struct {
	commonResult
}

// ScaleInResult is the response of a ScaleIn operations.
type ScaleInResult struct {
	commonResult
}

// ClusterPolicyPage contains a single page of all policies from a ListDetails call.
type ClusterPolicyPage struct {
	pagination.SinglePageBase
}

// GetPolicyResult is the response of a Get operations.
type GetPolicyResult struct {
	commonResult
}

type Cluster struct {
	Config          map[string]interface{} `json:"config"`
	CreatedAt       time.Time              `json:"-"`
	Data            map[string]interface{} `json:"data"`
	Dependents      map[string]interface{} `json:"dependents"`
	DesiredCapacity int                    `json:"desired_capacity"`
	Domain          string                 `json:"domain"`
	ID              string                 `json:"id"`
	InitAt          time.Time              `json:"-"`
	MaxSize         int                    `json:"max_size"`
	Metadata        map[string]interface{} `json:"metadata"`
	MinSize         int                    `json:"min_size"`
	Name            string                 `json:"name"`
	Nodes           []string               `json:"nodes"`
	Policies        []string               `json:"policies"`
	ProfileID       string                 `json:"profile_id"`
	ProfileName     string                 `json:"profile_name"`
	Project         string                 `json:"project"`
	Status          string                 `json:"status"`
	StatusReason    string                 `json:"status_reason"`
	Timeout         int                    `json:"timeout"`
	UpdatedAt       time.Time              `json:"-"`
	User            string                 `json:"user"`
}

func (r commonResult) Extract() (*Cluster, error) {
	var s struct {
		Cluster *Cluster `json:"cluster"`
	}

	err := r.ExtractInto(&s)
	return s.Cluster, err
}

type Action struct {
	Action string `json:"action"`
}

func (r ResizeResult) Extract() (string, error) {
	var s Action
	err := r.ExtractInto(&s)
	if err != nil {
		return s.Action, err
	}

	return s.Action, nil
}

func (r ScaleInResult) Extract() (string, error) {
	var s struct {
		Action string `json:"action"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}

	return s.Action, nil
}

// ClusterPage contains a single page of all clusters from a List call.
type ClusterPage struct {
	pagination.LinkedPageBase
}

func (page ClusterPage) IsEmpty() (bool, error) {
	clusters, err := ExtractClusters(page)
	return len(clusters) == 0, err
}

// ExtractCluster provides access to the list of clusters in a page acquired from the List operation.
func ExtractClusters(r pagination.Page) ([]Cluster, error) {
	var s struct {
		Clusters []Cluster `json:"clusters"`
	}
	err := (r.(ClusterPage)).ExtractInto(&s)
	return s.Clusters, err
}

func (r *Cluster) UnmarshalJSON(b []byte) error {
	type tmp Cluster
	var s struct {
		tmp
		CreatedAt string `json:"created_at"`
		InitAt    string `json:"init_at"`
		UpdatedAt string `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Cluster(s.tmp)

	if s.CreatedAt != "" {
		r.CreatedAt, err = time.Parse(gophercloud.RFC3339Milli, s.CreatedAt)
		if err != nil {
			return err
		}
	}

	if s.InitAt != "" {
		r.InitAt, err = time.Parse(gophercloud.RFC3339Milli, s.InitAt)
		if err != nil {
			return err
		}
	}

	if s.UpdatedAt != "" {
		r.UpdatedAt, err = time.Parse(gophercloud.RFC3339Milli, s.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}

type ClusterPolicy struct {
	ClusterID   string `json:"cluster_id"`
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

// IsEmpty determines if ClusterPolicyPage contains any results.
func (page ClusterPolicyPage) IsEmpty() (bool, error) {
	clusterPolicies, err := ExtractClusterPolicies(page)
	return len(clusterPolicies) == 0, err
}

// Extract provides access to the individual Policy returned by the Get and
// Create functions.
func (r GetPolicyResult) Extract() (*ClusterPolicy, error) {
	var s struct {
		ClusterPolicy *ClusterPolicy `json:"cluster_policy"`
	}
	err := r.ExtractInto(&s)
	return s.ClusterPolicy, err
}
