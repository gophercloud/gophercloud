package nodes

import (
	"encoding/json"
	"fmt"
	"time"

	"strings"

	"reflect"

	"github.com/gophercloud/gophercloud"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

func (r commonResult) ExtractActionFromLocation() (string, error) {
	if len(r.Header) > 0 {
		location := r.Header.Get("Location")
		actionID := strings.Split(location, "actions/")
		if len(actionID) >= 2 {
			return actionID[1], nil
		} else {
			return "", fmt.Errorf("ExtractAction failed. actionID=%s Location=%s", actionID, location)
		}
	}

	return "", fmt.Errorf("ExtractAction failed. Header=%s", r.Header)
}

// Extract provides access to the individual node returned by Get and extracts Node
func (r commonResult) Extract() (*Node, error) {
	var s struct {
		Node *Node `json:"node"`
	}
	err := r.ExtractInto(&s)
	if err != nil {
		return &Node{}, err
	}

	action, errAction := r.ExtractActionFromLocation()
	if errAction == nil {
		s.Node.ActionID = action
	}

	return s.Node, err
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
