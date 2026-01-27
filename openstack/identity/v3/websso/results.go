package websso

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
)

// GetResult is the response from a Get token metadata request.
type GetResult struct {
	gophercloud.Result
}

// ExtractToken interprets a GetResult as a Token.
func (r GetResult) ExtractToken() (*tokens.Token, error) {
	var s tokens.Token
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	s.ID = r.Header.Get("X-Subject-Token")
	return &s, nil
}

// ExtractTokenID returns the token ID from the result.
func (r GetResult) ExtractTokenID() (string, error) {
	return r.Header.Get("X-Subject-Token"), r.Err
}

// ExtractServiceCatalog returns the ServiceCatalog from the token response.
func (r GetResult) ExtractServiceCatalog() (*tokens.ServiceCatalog, error) {
	var s tokens.ServiceCatalog
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto extracts the token response into the provided struct.
func (r GetResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "token")
}

// CreateResult is the response from creating a scoped token.
type CreateResult struct {
	gophercloud.Result
}

// ExtractToken interprets a CreateResult as a Token.
func (r CreateResult) ExtractToken() (*tokens.Token, error) {
	var s tokens.Token
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	s.ID = r.Header.Get("X-Subject-Token")
	return &s, nil
}

// ExtractTokenID returns the token ID from the result.
func (r CreateResult) ExtractTokenID() (string, error) {
	return r.Header.Get("X-Subject-Token"), r.Err
}

// ExtractServiceCatalog returns the ServiceCatalog from the token response.
func (r CreateResult) ExtractServiceCatalog() (*tokens.ServiceCatalog, error) {
	var s tokens.ServiceCatalog
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto extracts the token response into the provided struct.
func (r CreateResult) ExtractInto(v any) error {
	return r.ExtractIntoStructPtr(v, "token")
}
