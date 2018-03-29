package policytypes

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// PolicyType represents a clustering policy type in the Openstack cloud
type PolicyType struct {
	Name          string                         `json:"name"`
	Version       string                         `json:"version"`
	SupportStatus map[string][]SupportStatusType `json:"support_status"`
}

// SupportStatusType represents the support status information for a clustering policy type
type SupportStatusType struct {
	Status string `json:"status"`
	Since  string `json:"since"`
}

// ExtractPolicyTypes interprets a page of results as a slice of PolicyTypes.
func ExtractPolicyTypes(r pagination.Page) ([]PolicyType, error) {
	var s struct {
		PolicyTypes []PolicyType `json:"policy_types"`
	}
	err := (r.(PolicyTypePage)).ExtractInto(&s)
	return s.PolicyTypes, err
}

// PolicyTypePage contains a single page of all policy types from a List call.
type PolicyTypePage struct {
	pagination.SinglePageBase
}

// IsEmpty determines if a PolicyType contains any results.
func (page PolicyTypePage) IsEmpty() (bool, error) {
	policyTypes, err := ExtractPolicyTypes(page)
	return len(policyTypes) == 0, err
}

// PolicyTypeDetail represents the detailed policy type information for a clustering policy type
type PolicyTypeDetail struct {
	Name          string                         `json:"name"`
	Schema        SchemaType                     `json:"schema"`
	SupportStatus map[string][]SupportStatusType `json:"support_status,omitempty"`
}

// SchemaType represents the schema of a clustering policy type
type SchemaType struct {
	// senlin.policy.affinity:
	AvailabilityZone   map[string]interface{} `json:"availability_zone,omitempty"`
	EnableDrsExtension map[string]interface{} `json:"enable_drs_extension,omitempty"`
	Servergroup        map[string]interface{} `json:"servergroup,omitempty"`

	// senlin.policy.batch:
	MaxBatchSize map[string]interface{} `json:"max_batch_size,omitempty"`
	MinInService map[string]interface{} `json:"min_in_service,omitempty"`
	PauseTime    map[string]interface{} `json:"pause_time,omitempty"`

	// senlin.policy.health:
	Detection map[string]interface{} `json:"detection,omitempty"`
	Recovery  map[string]interface{} `json:"recovery,omitempty"`

	// senlin.policy.scaling:
	Adjustment map[string]interface{} `json:"adjustment,omitempty"`
	Event      map[string]interface{} `json:"event,omitempty"`

	// senlin.policy.region_placement:
	Regions map[string]interface{} `json:"regions,omitempty"`

	// senlin.policy.loadbalance:
	HealthMonitor   map[string]interface{} `json:"health_monitor,omitempty"`
	LbStatusTimeout map[string]interface{} `json:"lb_status_timeout,omitempty"`
	Loadbalancer    map[string]interface{} `json:"loadbalancer,omitempty"`
	Pool            map[string]interface{} `json:"pool,omitempty"`
	Vip             map[string]interface{} `json:"vip,omitempty"`

	// senlin.policy.deletion:
	Criteria              map[string]interface{} `json:"criteria,omitempty"`
	DestroyAfterDeletion  map[string]interface{} `json:"destroy_after_deletion,omitempty"`
	GracePeriod           map[string]interface{} `json:"grace_period,omitempty"`
	Hooks                 map[string]interface{} `json:"hooks,omitempty"`
	ReduceDesiredCapacity map[string]interface{} `json:"reduce_desired_capacity,omitempty"`

	// senlin.policy.zone_placement:
	Zones map[string]interface{} `json:"zones,omitempty"`
}

// Extract provides access to the individual policy type returned by Get and extracts PolicyTypeDetail
func (r policyTypeResult) Extract() (*PolicyTypeDetail, error) {
	var s struct {
		PolicyType *PolicyTypeDetail `json:"policy_type"`
	}
	err := r.ExtractInto(&s)
	return s.PolicyType, err
}

type policyTypeResult struct {
	gophercloud.Result
}

type GetResult struct {
	policyTypeResult
}
