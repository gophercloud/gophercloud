package qos

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// QoS contains all the information associated with an OpenStack QoS specification.
type QoS struct {
	// Name is the name of the QoS.
	Name string `json:"name"`
	// Unique identifier for the QoS.
	ID string `json:"id"`
	// Consumer of QoS
	Consumer string `json:"consumer"`
	// Arbitrary key-value pairs defined by the user.
	Specs map[string]string `json:"specs"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the QoS object out of the commonResult object.
func (r commonResult) Extract() (*QoS, error) {
	var s QoS
	err := r.ExtractInto(&s)
	return &s, err
}

// ExtractInto converts our response data into a QoS struct
func (r commonResult) ExtractInto(qos interface{}) error {
	return r.Result.ExtractIntoStructPtr(qos, "qos_specs")
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	gophercloud.ErrResult
}

type QoSPage struct {
	pagination.LinkedPageBase
}

// IsEmpty determines if a QoSPage contains any results.
func (page QoSPage) IsEmpty() (bool, error) {
	qos, err := ExtractQoS(page)
	return len(qos) == 0, err
}

// NextPageURL uses the response's embedded link reference to navigate to the
// next page of results.
func (page QoSPage) NextPageURL() (string, error) {
	var s struct {
		Links []gophercloud.Link `json:"qos_specs_links"`
	}
	err := page.ExtractInto(&s)
	if err != nil {
		return "", err
	}
	return gophercloud.ExtractNextURL(s.Links)
}

// ExtractQoS provides access to the list of qos in a page acquired
// from the List operation.
func ExtractQoS(r pagination.Page) ([]QoS, error) {
	var s struct {
		QoSs []QoS `json:"qos_specs"`
	}
	err := (r.(QoSPage)).ExtractInto(&s)
	return s.QoSs, err
}

// GetResult is the response of a Get operations. Call its Extract method to
// interpret it as a Flavor.
type GetResult struct {
	commonResult
}

// Extract interprets any updateResult as qosSpecs, if possible.
func (r updateResult) Extract() (map[string]string, error) {
	var s struct {
		QosSpecs map[string]string `json:"qos_specs"`
	}
	err := r.ExtractInto(&s)
	return s.QosSpecs, err
}

// updateResult contains the result of a call for (potentially) multiple
// key-value pairs. Call its Extract method to interpret it as a
// map[string]interface.
type updateResult struct {
	gophercloud.Result
}

// AssociateResult contains the response body and error from a Associate request.
type AssociateResult struct {
	gophercloud.ErrResult
}

// DisassociateResult contains the response body and error from a Disassociate request.
type DisassociateResult struct {
	gophercloud.ErrResult
}

// DisassociateAllResult contains the response body and error from a DisassociateAll request.
type DisassociateAllResult struct {
	gophercloud.ErrResult
}
