package policies

import (
	"errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Validation errors returned by create or update operations.
var (
	ErrNoName = errors.New("Policy name cannot by empty.")
	ErrNoArgs = errors.New("Args cannot be nil for schedule policies.")
)

// List returns all scaling policies for a group.
func List(client *gophercloud.ServiceClient, groupID string) pagination.Pager {
	url := listURL(client, groupID)

	createPageFn := func(r pagination.PageResult) pagination.Page {
		return PolicyPage{pagination.SinglePageBase(r)}
	}

	return pagination.NewPager(client, url, createPageFn)
}

// CreateOptsBuilder is the interface responsible for generating the map that
// will be marshalled to JSON for a Create operation.
type CreateOptsBuilder interface {
	ToPolicyCreateMap() ([]map[string]interface{}, error)
}

// Adjustment represents the change in capacity associated with a policy.
type Adjustment struct {
	// The type for this adjustment.
	Type AdjustmentType

	// The value of the adjustment.  For adjustments of type Change or
	// DesiredCapacity, this will be converted to an integer.
	Value float64
}

// AdjustmentType represents the way in which a policy will change a group.
type AdjustmentType string

// Valid types of adjustments for a policy.
const (
	Change          AdjustmentType = "change"
	ChangePercent   AdjustmentType = "changePercent"
	DesiredCapacity AdjustmentType = "desiredCapacity"
)

// CreateOpts is a slice of CreateOpt structs that allow the user to create
// multiple policies in a single operation.
type CreateOpts []CreateOpt

// CreateOpt represents the options to create a policy.
type CreateOpt struct {
	// Name [required] is a name for the policy.
	Name string

	// Type [required] of policy, i.e. either "webhook" or "schedule".
	Type Type

	// Cooldown [required] period in seconds.
	Cooldown int

	// Adjustment [requried] type and value for the policy.
	Adjustment Adjustment

	// Additional configuration options for some types of policy.
	Args map[string]interface{}
}

// ToPolicyCreateMap converts a slice of CreateOpt structs into a map for use
// in the request body of a Create operation.
func (opts CreateOpts) ToPolicyCreateMap() ([]map[string]interface{}, error) {
	var policies []map[string]interface{}

	for _, o := range opts {
		if o.Name == "" {
			return nil, ErrNoName
		}

		if o.Type == Schedule && o.Args == nil {
			return nil, ErrNoArgs
		}

		policy := make(map[string]interface{})

		policy["name"] = o.Name
		policy["type"] = o.Type
		policy["cooldown"] = o.Cooldown

		// TODO: Function to validate and cast key + value?
		policy[string(o.Adjustment.Type)] = o.Adjustment.Value

		if o.Args != nil {
			policy["args"] = o.Args
		}

		policies = append(policies, policy)
	}

	return policies, nil
}

// Create requests a new policy be created and associated with the given group.
func Create(client *gophercloud.ServiceClient, groupID string, opts CreateOptsBuilder) CreateResult {
	var res CreateResult

	reqBody, err := opts.ToPolicyCreateMap()

	if err != nil {
		res.Err = err
		return res
	}

	_, res.Err = client.Post(createURL(client, groupID), reqBody, &res.Body, nil)

	return res
}

// Get requests the details of a single policy with the given ID.
func Get(client *gophercloud.ServiceClient, groupID, policyID string) GetResult {
	var result GetResult

	_, result.Err = client.Get(getURL(client, groupID, policyID), &result.Body, nil)

	return result
}
