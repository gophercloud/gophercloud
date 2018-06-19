package receivers

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

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
