package services

import (
	"time"

	"github.com/gophercloud/gophercloud/pagination"
)

type Service struct {
	Binary         string    `json:"binary"`
	CreatedAt      time.Time `json:"created_at"`
	State          string    `json:"state"`
	ReportCount    int       `json:"report_count"`
	UpdatedAt      time.Time `json:"updated_at"`
	Host           string    `json:"host"`
	Disabled       bool      `json:"disabled"`
	DisabledReason string    `json:"disabled_reason"`
	ID             int       `json:"id"`
}

type ServicePage struct {
	pagination.SinglePageBase
}

func (page ServicePage) IsEmpty() (bool, error) {
	services, err := ExtractServices(page)
	return len(services) == 0, err
}

func ExtractServices(r pagination.Page) ([]Service, error) {
	var s struct {
		Service []Service `json:"mservices"`
	}
	err := (r.(ServicePage)).ExtractInto(&s)
	return s.Service, err
}
