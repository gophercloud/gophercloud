package nodes

import (
	"encoding/json"
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

// Extract provides access to the individual node returned by Get and Create
func (r commonResult) Extract() (*Node, error) {
	var s struct {
		Node *Node `json:"node"`
	}
	err := r.ExtractInto(&s)

	return s.Node, err
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// ExtractNodes provides access to the list of nodes in a page acquired from the ListDetail operation.
func ExtractNodes(r pagination.Page) ([]Node, error) {
	var s struct {
		Nodes []Node `json:"nodes"`
	}
	err := (r.(NodePage)).ExtractInto(&s)
	return s.Nodes, err
}

// NodePage contains a single page of all nodes from a ListDetails call.
type NodePage struct {
	pagination.LinkedPageBase
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

// IsEmpty determines if a NodePage contains any results.
func (page NodePage) IsEmpty() (bool, error) {
	nodes, err := ExtractNodes(page)
	return len(nodes) == 0, err
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
