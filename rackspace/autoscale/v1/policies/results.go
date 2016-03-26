package policies

import (
	"github.com/mitchellh/mapstructure"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/pagination"
)

type policyResult struct {
	gophercloud.Result
}

// Policy represents a scaling policy.
type Policy struct {
	// UUID for the policy.
	ID string `mapstructure:"id" json:"id"`

	// Name of the policy.
	Name string `mapstructure:"name" json:"name"`

	// Type of scaling policy.
	Type Type `mapstructure:"type" json:"type"`

	// Cooldown period, in seconds.
	Cooldown int `mapstructure:"cooldown" json:"cooldown"`

	// Number of servers added or, if negative, removed.
	Change interface{} `mapstructure:"change" json:"change"`

	// Percent change to make in the number of servers.
	ChangePercent interface{} `mapstructure:"changePercent" json:"changePercent"`

	// Desired capacity of the of the associated group.
	DesiredCapacity interface{} `mapstructure:"desiredCapacity" json:"desiredCapacity"`

	// Additional configuration options for some types of policy.
	Args map[string]interface{} `mapstructure:"args" json:"args"`
}

// Type represents a type of scaling policy.
type Type string

const (
	// Schedule policies run at given times.
	Schedule Type = "schedule"

	// Webhook policies are triggered by HTTP requests.
	Webhook Type = "webhook"
)

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
	casted := page.(PolicyPage).Body

	var response struct {
		Policies []Policy `mapstructure:"policies"`
	}

	err := mapstructure.Decode(casted, &response)

	if err != nil {
		return nil, err
	}

	return response.Policies, err
}
