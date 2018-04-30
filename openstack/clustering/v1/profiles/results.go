package profiles

import (
	"encoding/json"
	"fmt"
	"time"

	"reflect"

	"github.com/gophercloud/gophercloud"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response of a Create operation.
type CreateResult struct {
	commonResult
}

// Extract provides access to Profile returned by the Get and Create functions.
func (r commonResult) Extract() (*Profile, error) {
	var s struct {
		Profile *Profile `json:"profile"`
	}
	err := r.ExtractInto(&s)
	return s.Profile, err
}

type Spec struct {
	Type       string                 `json:"type"`
	Version    string                 `json:"version"`
	Properties map[string]interface{} `json:"properties"`
}

// Profile represent a detailed profile
type Profile struct {
	CreatedAt time.Time              `json:"-"`
	Domain    string                 `json:"domain"`
	ID        string                 `json:"id"`
	Metadata  map[string]interface{} `json:"metadata"`
	Name      string                 `json:"name"`
	Project   string                 `json:"project"`
	Spec      Spec                   `json:"spec"`
	Type      string                 `json:"type"`
	UpdatedAt time.Time              `json:"-"`
	User      string                 `json:"user"`
}

func (r *Profile) UnmarshalJSON(b []byte) error {
	type tmp Profile
	var s struct {
		tmp
		CreatedAt interface{} `json:"created_at"`
		UpdatedAt interface{} `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Profile(s.tmp)

	switch t := s.CreatedAt.(type) {
	case string:
		if t != "" {
			r.CreatedAt, err = time.Parse(gophercloud.RFC3339Milli, t)
			if err != nil {
				return err
			}
		}
	case nil:
		r.CreatedAt = time.Time{}
	default:
		return fmt.Errorf("Invalid type for time. type=%v", reflect.TypeOf(s.CreatedAt))
	}

	switch t := s.UpdatedAt.(type) {
	case string:
		if t != "" {
			r.UpdatedAt, err = time.Parse(gophercloud.RFC3339Milli, t)
			if err != nil {
				return err
			}
		}
	case nil:
		r.UpdatedAt = time.Time{}
	default:
		return fmt.Errorf("Invalid type for time. type=%v", reflect.TypeOf(s.UpdatedAt))
	}

	return nil
}
