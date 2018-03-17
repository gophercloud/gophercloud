package clusters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
	"strings"
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

// ClusterPage contains a single page of all clusters from a ListDetails call.
type ClusterPage struct {
	pagination.LinkedPageBase
}

type Cluster struct {
	ActionID        string                 `json:"-"`
	Config          map[string]interface{} `json:"config"`
	CreatedAt       time.Time              `json:"-"`
	Data            map[string]interface{} `json:"data"`
	Dependents      map[string]interface{} `json:"dependents"`
	DesiredCapacity int                    `json:"desired_capacity"`
	DomainUUID      string                 `json:"domain"`
	ID              string                 `json:"id"`
	InitAt          time.Time              `json:"-"`
	MaxSize         int                    `json:"max_size"`
	Metadata        map[string]interface{} `json:"metadata"`
	MinSize         int                    `json:"min_size"`
	Name            string                 `json:"name"`
	Nodes           []string               `json:"nodes"`
	Policies        []string               `json:"policies"`
	ProfileID       string                 `json:"profile_id"`
	ProfileName     string                 `json:"profile_name"`
	ProjectUUID     string                 `json:"project"`
	Status          string                 `json:"status"`
	StatusReason    string                 `json:"status_reason"`
	Timeout         int                    `json:"timeout"`
	UpdatedAt       time.Time              `json:"-"`
	UserUUID        string                 `json:"user"`
}

func (r commonResult) ExtractCluster() (*Cluster, error) {
	var s struct {
		Cluster Cluster `json:"cluster"`
	}
	err := r.ExtractInto(&s)

	// Location: http://dev.senlin.cloud.blizzard.net:8778/v1/actions/625628cd-f877-44be-bde0-fec79f84e13d
	if err == nil && len(r.Header) > 0 {
		location := r.Header.Get("Location")
		actionID := strings.Split(location, "actions/")
		if len(actionID) >= 2 {
			s.Cluster.ActionID = actionID[1]
		}
	}

	return &s.Cluster, err
}

func (r commonResult) Extract() (*Cluster, error) {
	s, err := r.ExtractCluster()

	// Location: http://dev.senlin.cloud.blizzard.net:8778/v1/actions/625628cd-f877-44be-bde0-fec79f84e13d
	if err == nil && len(r.Header) > 0 {
		location := r.Header.Get("Location")
		actionID := strings.Split(location, "actions/")
		if len(actionID) >= 2 {
			s.ActionID = actionID[1]
		}
	}

	return s, err
}

// ExtractCluster provides access to the list of clusters in a page acquired from the ListDetail operation.
func ExtractClusters(r pagination.Page) ([]Cluster, error) {
	var s struct {
		Clusters []Cluster `json:"clusters"`
	}
	err := (r.(ClusterPage)).ExtractInto(&s)
	return s.Clusters, err
}

// IsEmpty determines if a ClusterPage contains any results.
func (page ClusterPage) IsEmpty() (bool, error) {
	clusters, err := ExtractClusters(page)
	return len(clusters) == 0, err
}

type ClusterResult struct {
	Cluster string `json:"clusters"`
}

// Extract provides access to the individual Cluster returned by the Get and
// Create functions.
func (r commonResult) ExtractClusters() (s *ClusterResult, err error) {
	err = r.ExtractInto(&s)
	return s, err
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

func (r *Cluster) UnmarshalJSON(b []byte) error {
	type tmp Cluster
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
	*r = Cluster(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.InitAt = time.Time(s.InitAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}
