package policies

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// PostResult is the response of a Post operations.
type PostResult struct {
	commonResult
}

// DeleteResult is the result from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// UpdateResult is the response of a Update operations.
type UpdateResult struct {
	commonResult
}

// Extract provides access to the individual Flavor returned by the Get and
// Create functions.
func (r commonResult) Extract() (*Policy, error) {
	var s struct {
		Policy *Policy `json:"policy"`
	}
	err := r.ExtractInto(&s)
	return s.Policy, err
}

// Policy represents a detailed policy
type Policy struct {
	CreatedAt time.Time `json:"-"`

	Data map[string]interface{} `json:"data"`

	DomainUUID string `json:"domain"`

	ID string `json:"id"`

	Name string `json:"name"`

	ProjectUUID string `json:"project"`

	Spec map[string]interface{} `json:"spec"`

	Type string `json:"type"`

	UpdatedAt time.Time `json:"-"`

	UserUUID string `json:"user"`
}

// ExtractPolicies provides access to the list of profiles in a page acquired from the ListDetail operation.
func ExtractPolicies(r pagination.Page) ([]Policy, error) {
	var s struct {
		Policies []Policy `json:"policies"`
	}
	err := (r.(PolicyPage)).ExtractInto(&s)
	return s.Policies, err
}

// PolicyPage contains a single page of all policies from a ListDetails call.
type PolicyPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines if a ProfilePage contains any results.
func (page PolicyPage) IsEmpty() (bool, error) {
	policies, err := ExtractPolicies(page)
	return len(policies) == 0, err
}

// TODO: Should consolidate this. Should time parsing be strict
type JSONRFC3339Milli time.Time

func (jt *JSONRFC3339Milli) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*jt = JSONRFC3339Milli(time.Time{})
		return nil
	}

	b := bytes.NewBuffer(data)
	dec := json.NewDecoder(b)
	var s string
	if err := dec.Decode(&s); err != nil {
		return err
	}
	t, err := time.Parse(gophercloud.RFC3339Milli, s)
	if err != nil {
		t, err = time.Parse(gophercloud.RFC3339MilliNoZ, s)
		if err != nil {
			t, err = time.Parse(gophercloud.RFC3339NoZ, s)
			if err != nil {
				return err
			}
		}
	}
	*jt = JSONRFC3339Milli(t)
	return nil
}

func (r *Policy) UnmarshalJSON(b []byte) error {
	type tmp Policy
	var s struct {
		tmp
		CreatedAt JSONRFC3339Milli `json:"created_at"`
		UpdatedAt JSONRFC3339Milli `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error Unmarshal")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}
	*r = Policy(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}
