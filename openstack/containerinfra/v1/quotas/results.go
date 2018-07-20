package quotas

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

// CreateResult is the response of a Create operations.
type CreateResult struct {
	commonResult
}

// Extract is a function that accepts a result and extracts a quota resource.
func (r commonResult) Extract() (*Quotas, error) {
	var s *Quotas
	err := r.ExtractInto(&s)
	return s, err
}

type Quotas struct {
	Resource  string    `json:"resource"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	HardLimit int       `json:"hard_limit"`
	ProjectID string    `json:"project_id"`
	ID        string    `json:"-"`
}

func (r *Quotas) UnmarshalJSON(b []byte) error {
	type tmp Quotas
	var s struct {
		tmp
		ID interface{} `json:"id"`
	}

	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Quotas(s.tmp)

	switch t := s.ID.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		r.ID = fmt.Sprint(t)
	case string:
		r.ID = t
	}

	return nil
}
