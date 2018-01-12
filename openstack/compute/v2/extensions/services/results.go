package services

import "github.com/gophercloud/gophercloud/pagination"

// Service represents a Compute service in the OpenStack cloud.
type Service struct {
	// The binary name of the service.
	Binary string `json:"binary"`

	// The reason for disabling a service.
	DisabledReason string `json:"disabled_reason"`

	// Whether or not this service was forced down manually by an administrator.
	ForcedDown bool `json:"forced_down"`

	// The name of the host.
	Host string `json:"host"`

	// The id of the service.
	ID int `json:"id"`

	// The state of the service. One of up or down.
	State string `json:"state"`

	// The status of the service. One of enabled or disabled.
	Status string `json:"status"`

	// The availability zone name.
	Zone string `json:"zone"`
}

// ServicePage represents a single page of all Services from a List request.
type ServicePage struct {
	pagination.SinglePageBase
}

// IsEmpty determines whether or not a page of Services contains any results.
func (page ServicePage) IsEmpty() (bool, error) {
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
