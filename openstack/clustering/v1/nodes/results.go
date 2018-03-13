package nodes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"strings"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response of a Create operations.
type CreateResult struct {
	commonResult
}

// DeleteResult is the result from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// PostResult is the response of a Post operations.
type PostResult struct {
	commonResult
}

// UpdateResult is the response of a Update operations.
type UpdateResult struct {
	commonResult
}

// Extract provides access to the individual node returned by Get and extracts Node
func (r commonResult) Extract() (*Node, error) {
	var s struct {
		Node *Node `json:"node"`
	}
	err := r.ExtractInto(&s)

	// Location: http://dev.senlin.cloud.blizzard.net:8778/v1/actions/625628cd-f877-44be-bde0-fec79f84e13d
	if err == nil && len(r.Header) > 0 {
		location := r.Header.Get("Location")
		actionID := strings.Split(location, "actions/")
		if len(actionID) >= 2 {
			s.Node.ActionID = actionID[1]
		}
	}

	return s.Node, err
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
	ActionID string `json:"-"`

	ClusterUUID string `json:"cluster_id"`

	CreatedAt time.Time `json:"-"`

	Data DataType `json:"data"`

	Dependents map[string]interface{} `json:"dependents"`

	DomainUUID string `json:"domain"`

	ID string `json:"id"`

	Index int `json:"index"`

	InitAt time.Time `json:"-"`

	Metadata map[string]interface{} `json:"metadata"`

	Name string `json:"name"`

	PhysicalUUID string `json:"physical_id"`

	ProfileUUID string `json:"profile_id"`

	ProfileName string `json:"profile_name"`

	ProjectUUID string `json:"project_id"`

	Role string `json:"role"`

	Status string `json:"status"`

	StatusReason string `json:"status_reason"`

	UpdatedAt time.Time `json:"-"`

	UserUUID string `json:"user"`
}

// IsEmpty determines if a NodePage contains any results.
func (page NodePage) IsEmpty() (bool, error) {
	nodes, err := ExtractNodes(page)
	return len(nodes) == 0, err
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

func (r *Node) UnmarshalJSON(b []byte) error {
	type tmp Node
	var s struct {
		tmp
		CreatedAt JSONRFC3339Milli `json:"created_at"`
		InitAt    JSONRFC3339Milli `json:"init_at"`
		UpdatedAt JSONRFC3339Milli `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		fmt.Println("Error Unmarshal")
		fmt.Printf("Detail Unmarshal Error: %v", err)
		return err
	}
	*r = Node(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.InitAt = time.Time(s.InitAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}
