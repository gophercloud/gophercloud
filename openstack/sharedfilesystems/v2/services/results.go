package services

import (
	"encoding/json"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Service represents a Shared File System service in the OpenStack cloud.
type Service struct {
	// The binary name of the service.
	Binary string `json:"binary"`

	// The name of the host.
	Host string `json:"host"`

	// The ID of the service.
	ID int `json:"id"`

	// The state of the service. One of up or down.
	State string `json:"state"`

	// The status of the service. One of available or unavailable.
	Status string `json:"status"`

	// The date and time stamp when the extension was last updated.
	UpdatedAt time.Time `json:"-"`

	// The availability zone name.
	Zone string `json:"zone"`
}

// UnmarshalJSON to override default
func (r *Service) UnmarshalJSON(b []byte) error {
	type tmp Service
	var s struct {
		tmp
		UpdatedAt gophercloud.JSONRFC3339MilliNoZ `json:"updated_at"`
	}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*r = Service(s.tmp)

	r.UpdatedAt = time.Time(s.UpdatedAt)

	return nil
}

// ServicePage represents a single page of all Services from a List request.
type ServicePage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a page of Services contains any results.
func (page ServicePage) IsEmpty() (bool, error) {
	if page.StatusCode == 204 {
		return true, nil
	}

	services, err := ExtractServices(page)
	return len(services) == 0, err
}

func ExtractServices(r pagination.Page) ([]Service, error) {
	var s struct {
		Service []Service `json:"services"`
	}
	err := (r.(ServicePage)).ExtractInto(&s)
	return s.Service, err
}
