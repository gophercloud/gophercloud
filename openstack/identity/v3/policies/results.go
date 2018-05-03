package policies

import (
	"encoding/json"

	"github.com/gophercloud/gophercloud/internal"
	"github.com/gophercloud/gophercloud/pagination"
)

// Policy is an arbitrarily serialized policy engine rule
// set to be consumed by a remote service.
type Policy struct {
	// ID is the unique ID of the policy.
	ID string `json:"id"`

	// Blob is the policy rule as a serialized blob.
	Blob string `json:"blob"`

	// Type is the MIME media type of the serialized policy blob.
	Type string `json:"type"`

	// Links contains referencing links to the policy.
	Links map[string]interface{} `json:"links"`

	// Extra is a collection of miscellaneous key/values.
	Extra map[string]interface{} `json:"-"`
}

func (r *Policy) UnmarshalJSON(b []byte) error {
	type tmp Policy
	var s struct {
		tmp
		Extra map[string]interface{} `json:"extra"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Policy(s.tmp)

	// Collect other fields and bundle them into Extra
	// but only if a field titled "extra" wasn't sent.
	if s.Extra != nil {
		r.Extra = s.Extra
	} else {
		var result interface{}
		err := json.Unmarshal(b, &result)
		if err != nil {
			return err
		}
		if resultMap, ok := result.(map[string]interface{}); ok {
			r.Extra = internal.RemainingKeys(Policy{}, resultMap)
		}
	}

	return err
}

// PolicyPage is a single page of Policy results.
type PolicyPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines whether or not a page of Policies contains any results.
func (r PolicyPage) IsEmpty() (bool, error) {
	policies, err := ExtractPolicies(r)
	return len(policies) == 0, err
}

// NextPageURL extracts the "next" link from the links section of the result.
func (r PolicyPage) NextPageURL() (string, error) {
	var s struct {
		Links struct {
			Next     string `json:"next"`
			Previous string `json:"previous"`
		} `json:"links"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return s.Links.Next, err
}

// ExtractPolicies returns a slice of Policies
// contained in a single page of results.
func ExtractPolicies(r pagination.Page) ([]Policy, error) {
	var s struct {
		Policies []Policy `json:"policies"`
	}
	err := (r.(PolicyPage)).ExtractInto(&s)
	return s.Policies, err
}
