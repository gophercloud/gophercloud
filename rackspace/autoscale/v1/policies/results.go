package policies

import (
	"github.com/mitchellh/mapstructure"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type policyResult struct {
	gophercloud.Result
}

// Extract interprets any policyResult as a Policy, if possible.
func (r policyResult) Extract() (*Policy, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Policy policy `mapstructure:"policy"`
	}

	if err := mapstructure.Decode(r.Body, &response); err != nil {
		return nil, err
	}

	policy := response.Policy.toExported()

	return &policy, nil
}

// CreateResult represents the result of a create operation.
type CreateResult struct {
	policyResult
}

// Extract extracts a slice of Policies from a CreateResult.  Multiple policies
// can be created in a single operation, so the result of a create is always a
// list of policies.
func (res CreateResult) Extract() ([]Policy, error) {
	if res.Err != nil {
		return nil, res.Err
	}

	return commonExtractPolicies(res.Body)
}

// GetResult temporarily contains the response from a Get call.
type GetResult struct {
	policyResult
}

// UpdateResult represents the result of an update operation.
type UpdateResult struct {
	gophercloud.ErrResult
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}

// ExecuteResult represents the result of an execute operation.
type ExecuteResult struct {
	gophercloud.ErrResult
}

// Type represents a type of scaling policy.
type Type string

const (
	// Schedule policies run at given times.
	Schedule Type = "schedule"

	// Webhook policies are triggered by HTTP requests.
	Webhook Type = "webhook"
)

// AdjustmentType represents the way in which a policy will change a group.
type AdjustmentType string

// Valid types of adjustments for a policy.
const (
	Change          AdjustmentType = "change"
	ChangePercent   AdjustmentType = "changePercent"
	DesiredCapacity AdjustmentType = "desiredCapacity"
)

// Policy represents a scaling policy.
type Policy struct {
	// UUID for the policy.
	ID string

	// Name of the policy.
	Name string

	// Type of scaling policy.
	Type Type

	// Cooldown period, in seconds.
	Cooldown int

	// The type of adjustment in capacity to be made.
	AdjustmentType AdjustmentType

	// The numeric value of the adjustment in capacity.
	AdjustmentValue float64

	// Additional configuration options for some types of policy.
	Args map[string]interface{}
}

// This is an intermediate representation of the exported Policy type.  The
// fields in API responses vary by policy type and configuration.  This lets us
// decode responses then normalize them into a Policy.
type policy struct {
	ID       string `mapstructure:"id"`
	Name     string `mapstructure:"name"`
	Type     Type   `mapstructure:"type"`
	Cooldown int    `mapstructure:"cooldown"`

	// The API will respond with exactly one of these omitting the others.
	Change          interface{} `mapstructure:"change"`
	ChangePercent   interface{} `mapstructure:"changePercent"`
	DesiredCapacity interface{} `mapstructure:"desiredCapacity"`

	// Additional configuration options for schedule policies.
	Args map[string]interface{} `mapstructure:"args"`
}

// Assemble a Policy from the intermediate policy struct.
func (p policy) toExported() Policy {
	policy := Policy{}

	policy.ID = p.ID
	policy.Name = p.Name
	policy.Type = p.Type
	policy.Cooldown = p.Cooldown

	policy.Args = p.Args

	if v, ok := p.Change.(float64); ok {
		policy.AdjustmentType = Change
		policy.AdjustmentValue = v
	} else if v, ok := p.ChangePercent.(float64); ok {
		policy.AdjustmentType = ChangePercent
		policy.AdjustmentValue = v
	} else if v, ok := p.DesiredCapacity.(float64); ok {
		policy.AdjustmentType = DesiredCapacity
		policy.AdjustmentValue = v
	}

	return policy
}

// PolicyPage is the page returned by a pager when traversing over a collection
// of scaling policies.
type PolicyPage struct {
	pagination.SinglePageBase
}

// IsEmpty returns true if a page contains no Policy results.
func (page PolicyPage) IsEmpty() (bool, error) {
	policies, err := ExtractPolicies(page)

	if err != nil {
		return true, err
	}

	return len(policies) == 0, nil
}

// ExtractPolicies interprets the results of a single page from a List() call,
// producing a slice of Policies.
func ExtractPolicies(page pagination.Page) ([]Policy, error) {
	return commonExtractPolicies(page.(PolicyPage).Body)
}

func commonExtractPolicies(body interface{}) ([]Policy, error) {
	var response struct {
		Policies []policy `mapstructure:"policies"`
	}

	err := mapstructure.Decode(body, &response)

	if err != nil {
		return nil, err
	}

	policies := make([]Policy, len(response.Policies))

	for i, p := range response.Policies {
		policies[i] = p.toExported()
	}

	return policies, nil
}
