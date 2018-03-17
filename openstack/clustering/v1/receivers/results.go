package receivers

import (
	"bytes"
	"encoding/json"
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
func (r commonResult) Extract() (*Receiver, error) {
	var s struct {
		Receiver *Receiver `json:"receiver"`
	}
	err := r.ExtractInto(&s)
	return s.Receiver, err
}

// ExtractReceivers provides access to the list of nodes in a page acquired from the ListDetail operation.
func ExtractReceivers(r pagination.Page) ([]Receiver, error) {
	var s struct {
		Receivers []Receiver `json:"receivers"`
	}
	err := (r.(ReceiverPage)).ExtractInto(&s)
	return s.Receivers, err
}

// ReceiverPage contains a single page of all nodes from a ListDetails call.
type ReceiverPage struct {
	pagination.LinkedPageBase
}

// Receiver represent a detailed receiver
type Receiver struct {
	Action      string                 `json:"action,omitempty"`
	Actor       map[string]interface{} `json:"actor,omitempty"`
	Channel     map[string]interface{} `json:"channel,omitempty"`
	ClusterUUID string                 `json:"cluster_id"`
	CreatedAt   time.Time              `json:"-"`
	DomainUUID  string                 `json:"domain"`
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Params      map[string]interface{} `json:"params,omitempty"`
	ProjectUUID string                 `json:"project"`
	Type        string                 `json:"type"`
	UpdatedAt   time.Time              `json:"-"`
	UserUUID    string                 `json:"user"`
}

// IsEmpty determines if a NodePage contains any results.
func (page ReceiverPage) IsEmpty() (bool, error) {
	receivers, err := ExtractReceivers(page)
	return len(receivers) == 0, err
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

func (r *Receiver) UnmarshalJSON(b []byte) error {
	type tmp Receiver
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
	*r = Receiver(s.tmp)

	r.CreatedAt = time.Time(s.CreatedAt)
	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}
