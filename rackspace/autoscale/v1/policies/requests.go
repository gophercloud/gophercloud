package policies

import (
	"errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

// Validation errors returned by create or update operations.
var (
	ErrNoName            = errors.New("Policy name cannot by empty.")
	ErrNoArgs            = errors.New("Args cannot be nil for schedule policies.")
	ErrInvalidAdjustment = errors.New("Invalid adjustment type.")
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

	// AdjustmentType [requried] is the method used to change the capacity of
	// the group, i.e. one of: Change, ChangePercent, or DesiredCapacity.
	AdjustmentType AdjustmentType

	// AdjustmentValue [required] is the numeric value of the adjustment.  For
	// adjustments of type Change or DesiredCapacity, this will be converted to
	// an integer.
	AdjustmentValue float64

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

		err := setAdjustment(o.AdjustmentType, o.AdjustmentValue, policy)

		if err != nil {
			return nil, err
		}

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

// UpdateOptsBuilder is the interface responsible for generating the map
// structure for producing JSON for an Update operation.
type UpdateOptsBuilder interface {
	ToPolicyUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts represents the options for updating an existing policy.
//
// Update operations completely replace the configuration being updated. Empty
// values in the update are accepted and overwrite previously specified
// parameters.
type UpdateOpts struct {
	// Name [required] is a name for the policy.
	Name string

	// Type [required] of policy, i.e. either "webhook" or "schedule".
	Type Type

	// Cooldown [required] period in seconds.  If you don't specify a cooldown,
	// it will default to zero, and the policy will be configured as such.
	Cooldown int

	// AdjustmentType [requried] is the method used to change the capacity of
	// the group, i.e. one of: Change, ChangePercent, or DesiredCapacity.
	AdjustmentType AdjustmentType

	// AdjustmentValue [required] is the numeric value of the adjustment.  For
	// adjustments of type Change or DesiredCapacity, this will be converted to
	// an integer.
	AdjustmentValue float64

	// Additional configuration options for some types of policy.
	Args map[string]interface{}
}

// ToPolicyUpdateMap converts an UpdateOpts struct into a map for use as the
// request body in an Update request.
func (opts UpdateOpts) ToPolicyUpdateMap() (map[string]interface{}, error) {
	if opts.Name == "" {
		return nil, ErrNoName
	}

	if opts.Type == Schedule && opts.Args == nil {
		return nil, ErrNoArgs
	}

	policy := make(map[string]interface{})

	policy["name"] = opts.Name
	policy["type"] = opts.Type
	policy["cooldown"] = opts.Cooldown

	err := setAdjustment(opts.AdjustmentType, opts.AdjustmentValue, policy)

	if err != nil {
		return nil, err
	}

	if opts.Args != nil {
		policy["args"] = opts.Args
	}

	return policy, nil
}

// Update requests the configuration of the given policy be updated.
func Update(client *gophercloud.ServiceClient, groupID, policyID string, opts UpdateOptsBuilder) UpdateResult {
	var result UpdateResult

	url := updateURL(client, groupID, policyID)
	reqBody, err := opts.ToPolicyUpdateMap()

	if err != nil {
		result.Err = err
		return result
	}

	_, result.Err = client.Put(url, reqBody, nil, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})

	return result
}

// Delete requests the given policy be permanently deleted.
func Delete(client *gophercloud.ServiceClient, groupID, policyID string) DeleteResult {
	var result DeleteResult

	url := deleteURL(client, groupID, policyID)
	_, result.Err = client.Delete(url, &gophercloud.RequestOpts{
		OkCodes: []int{204},
	})

	return result
}

// Execute requests the given policy be executed immediately.
func Execute(client *gophercloud.ServiceClient, groupID, policyID string) ExecuteResult {
	var result ExecuteResult

	url := executeURL(client, groupID, policyID)
	_, result.Err = client.Post(url, nil, &result.Body, &gophercloud.RequestOpts{
		OkCodes: []int{202},
	})

	return result
}

// Validate and set an adjustment on the given request body.
func setAdjustment(t AdjustmentType, v float64, body map[string]interface{}) error {
	key := string(t)

	switch t {
	case ChangePercent:
		body[key] = v

	case Change, DesiredCapacity:
		body[key] = int(v)

	default:
		return ErrInvalidAdjustment
	}

	return nil
}
