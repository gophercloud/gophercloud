package nodes

import (
	"encoding/json"
	"fmt"
	"time"

	"reflect"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
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

type Placement struct {
	Zone string `json:"zone"`
}

type InternalPort struct {
	FixedIps         []FixedIp `json:"fixed_ips"`
	ID               string    `json:"id"`
	NetworkID        string    `json:"network_id"`
	Remove           bool      `json:"remove"`
	SecurityGroupIds []string  `json:"security_group_ids"`
}

type FixedIp struct {
	IPAddress string `json:"ip_address"`
	SubnetID  string `json:"subnet_id"`
}

type DataType struct {
	InternalPorts []InternalPort `json:"internal_ports"`
	Placement     Placement      `json:"placement"`
}

// Node represents a node structure
type Node struct {
	ActionID     string                 `json:"-"`
	ClusterID    string                 `json:"cluster_id"`
	CreatedAt    time.Time              `json:"-"`
	Data         DataType               `json:"data"`
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
