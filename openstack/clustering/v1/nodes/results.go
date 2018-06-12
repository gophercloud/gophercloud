package nodes

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/gophercloud/gophercloud"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response of a Create operations.
type CreateResult struct {
	commonResult
}

// Extract provides access to the individual node returned by Get and Create
func (r commonResult) Extract() (*Node, error) {
	var s struct {
		Node *Node `json:"node"`
	}
	err := r.ExtractInto(&s)

	return s.Node, err
}

// Node represents a node structure
type Node struct {
	ClusterID    string                 `json:"cluster_id"`
	CreatedAt    time.Time              `json:"-"`
	Data         map[string]interface{} `json:"data"`
	Dependents   map[string]interface{} `json:"dependents"`
	Domain       string                 `json:"domain"`
	ID           string                 `json:"id"`
	Index        int                    `json:"index"`
	InitAt       time.Time              `json:"-"`
	Metadata     map[string]interface{} `json:"metadata"`
	Name         string                 `json:"name"`
	PhysicalID   string                 `json:"physical_id"`
	ProfileID    string                 `json:"profile_id"`
	ProfileName  string                 `json:"profile_name"`
	ProjectID    string                 `json:"project_id"`
	Role         string                 `json:"role"`
	Status       string                 `json:"status"`
	StatusReason string                 `json:"status_reason"`
	UpdatedAt    time.Time              `json:"-"`
	User         string                 `json:"user"`
}

func (r *Node) UnmarshalJSON(b []byte) error {
	type tmp Node
	var s struct {
		tmp
		CreatedAt interface{} `json:"created_at"`
		InitAt    interface{} `json:"init_at"`
		UpdatedAt interface{} `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Node(s.tmp)

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

	switch t := s.InitAt.(type) {
	case string:
		if t != "" {
			r.InitAt, err = time.Parse(gophercloud.RFC3339Milli, t)
			if err != nil {
				return err
			}
		}
	case nil:
		r.InitAt = time.Time{}
	default:
		return fmt.Errorf("Invalid type for time. type=%v", reflect.TypeOf(s.InitAt))
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
