package nodes

import (
	"encoding/json"
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
		CreatedAt string `json:"created_at"`
		InitAt    string `json:"init_at"`
		UpdatedAt string `json:"updated_at"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Node(s.tmp)

	if s.CreatedAt != "" {
		r.CreatedAt, err = time.Parse(time.RFC3339, s.CreatedAt)
		if err != nil {
			return err
		}
	}

	if s.InitAt != "" {
		r.InitAt, err = time.Parse(time.RFC3339, s.InitAt)
		if err != nil {
			return err
		}
	}

	if s.UpdatedAt != "" {
		r.UpdatedAt, err = time.Parse(time.RFC3339, s.UpdatedAt)
		if err != nil {
			return err
		}
	}

	return nil
}
