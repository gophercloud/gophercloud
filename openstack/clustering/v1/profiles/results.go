package profiles

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"time"
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
func (r commonResult) Extract() (*Profile, error) {
	var s struct {
		Profile *Profile `json:"profile"`
	}
	err := r.ExtractInto(&s)
	return s.Profile, err
}

// Profile represent a detailed profile
type Profile struct {
	CreatedAt time.Time `json:"-"`

	DomainUUID string `json:"domain"`

	ID string `json:"id"`

	Metadata map[string]interface{} `json:"metadata"`

	Name string `json:"name"`

	ProjectUUID string `json:"project"`

	Spec map[string]interface{} `json:"spec"`

	Type string `json:"type"`

	UpdatedAt time.Time `json:"-"`

	UserUUID string `json:"user"`
}

// ExtractProfiles provides access to the list of profiles in a page acquired from the ListDetail operation.
func ExtractProfiles(r pagination.Page) ([]Profile, error) {
	var s struct {
		Profiles []Profile `json:"profiles"`
	}
	err := (r.(ProfilePage)).ExtractInto(&s)
	return s.Profiles, err
}

// NodePage contains a single page of all nodes from a ListDetails call.
type ProfilePage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines if a ProfilePage contains any results.
func (page ProfilePage) IsEmpty() (bool, error) {
	profiles, err := ExtractProfiles(page)
	return len(profiles) == 0, err
}

type JSONRFC3339Milli time.Time

const RFC3339Milli = "2006-01-02T15:04:05.999999Z"

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
	t, err := time.Parse(RFC3339Milli, s)
	if err != nil {
		return err
	}
	*jt = JSONRFC3339Milli(t)
	return nil
}

type ProfileSpec struct {
	profileSpec map[string]interface{} `json:"spec"`
}

func (ps *ProfileSpec) UnmarshalJSON(b []byte) (err error) {
	var m interface{}

	err = json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println("Error Unmarshal Profile Spec")
		fmt.Printf("Detail Unmarshal Profile Spec Error: %v", err)
		return err
	}

	// Profile Spec should only be map[string]interface{}
	switch typeValue := m.(type) {
	case map[string]interface{}:
		if nw, isGood := typeValue["properties"].(map[string]interface{})["networks"]; isGood {
			switch networks := nw.(type) {
			case []interface{}:
				for _, network := range networks {
					if false {
						fmt.Println("Profile-Spec-Network", network)
					}
				}
			default:
				err = errors.New("Unknown type for profile spec network. Only supports map[string]interface{}")
				fmt.Printf("Error Unmarshal Profile Spec Network Error: %v", err)
				return err
			}
		}
	default:
		err = errors.New("Unknown type for profile spec. Only supports map[string]interface{}")
		fmt.Printf("Error Unmarshal Profile Spec Error: %v", err)
		return err
	}

	*ps = ProfileSpec{profileSpec: m.(map[string]interface{})}
	return nil
}

func (r *Profile) UnmarshalJSON(b []byte) error {
	type tmp Profile
	var s struct {
		tmp
		CreatedAt JSONRFC3339Milli `json:"created_at"`
		UpdatedAt JSONRFC3339Milli `json:"updated_at"`
		Spec      ProfileSpec      `json:"spec"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error Unmarshal")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}
	*r = Profile(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)
	r.Spec = s.Spec.profileSpec

	return nil
}
