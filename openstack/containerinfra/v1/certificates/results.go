package certificates

import (
	"github.com/gophercloud/gophercloud"
)

type commonResult struct {
	gophercloud.Result
}

type createResult struct {
	gophercloud.Result
}

// GetResult is the response of a Get operations.
type GetResult struct {
	commonResult
}

// CreateResult is the response of a Create operations.
type CreateResult struct {
	createResult
}

// Extract is a function that accepts a result and extracts a certificate resource.
func (r commonResult) Extract() (*Certificate, error) {
	var s *Certificate
	err := r.ExtractInto(&s)
	return s, err
}

// Extract is a function that accepts a result and extracts a CreateCertificateResponsel resource.
func (r createResult) Extract() (*CreateCertificateResponse, error) {
	var s *CreateCertificateResponse
	err := r.ExtractInto(&s)
	return s, err
}

// Represents a Certificate
type Certificate struct {
	ClusterUUID string             `json:"cluster_uuid"`
	BayUUID     string             `json:"bay_uuid"`
	Links       []gophercloud.Link `json:"links"`
	Pem         string             `json:"pem"`
}

// Represents a Certificate Create Response
type CreateCertificateResponse struct {
	BayUUID     string             `json:"bay_uuid"`
	Links       []gophercloud.Link `json:"links"`
	Pem         string             `json:"pem"`
	Csr         string             `json:"csr"`
}

