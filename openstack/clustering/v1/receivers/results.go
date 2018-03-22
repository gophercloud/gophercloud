package receivers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// commonResult is the response of a base result.
type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response for a create operation.
type CreateResult struct {
	commonResult
}

// GetResult is the response for a get operation.
type GetResult struct {
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

// Receiver represent a detailed receiver
type Receiver struct {
	Action    string                 `json:"action"`
	Actor     map[string]interface{} `json:"actor"`
	Channel   map[string]interface{} `json:"channel"`
	Cluster   string                 `json:"cluster_id"`
	CreatedAt time.Time              `json:"-"`
	Domain    string                 `json:"domain"`
	ID        string                 `json:"id"`
	Name      string                 `json:"name"`
	Params    map[string]interface{} `json:"params"`
	Project   string                 `json:"project"`
	Type      string                 `json:"type"`
	UpdatedAt time.Time              `json:"-"`
	User      string                 `json:"user"`
}

func (r *Receiver) UnmarshalJSON(b []byte) error {
	type tmp Receiver
	var s struct {
		tmp
		CreatedAt interface{} `json:"created_at"`
		UpdatedAt interface{} `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Receiver(s.tmp)

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

// IsEmpty determines if a ProfilePage contains any results.
func (page ReceiverPage) IsEmpty() (bool, error) {
	receivers, err := ExtractReceivers(page)
	return len(receivers) == 0, err
}

// DeleteResult is the result from a Delete operation. Call its ExtractErr
// method to determine if the call succeeded or failed.
type DeleteResult struct {
	gophercloud.ErrResult
}
