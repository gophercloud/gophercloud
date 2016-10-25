package sharetypes

import "github.com/gophercloud/gophercloud"

// ShareType contains all the information associated with an OpenStack
// ShareType.
type ShareType struct {
	// The Share Type ID
	ID string `json:"id"`
	// The Share Type name
	Name string `json:"name"`
	// Indicates whether a share type is publicly accessible
	IsPublic bool `json:"os-share-type-access:is_public"`
	// The required extra specifications for the share type
	RequiredExtraSpecs map[string]interface{} `json:"required_extra_specs"`
	// The extra specifications for the share type
	ExtraSpecs map[string]interface{} `json:"extra_specs"`
}

type commonResult struct {
	gophercloud.Result
}

// Extract will get the ShareType object out of the commonResult object.
func (r commonResult) Extract() (*ShareType, error) {
	var s struct {
		ShareType *ShareType `json:"share_type"`
	}
	err := r.ExtractInto(&s)
	return s.ShareType, err
}

// CreateResult contains the response body and error from a Create request.
type CreateResult struct {
	commonResult
}

// DeleteResult contains the response body and error from a Delete request.
type DeleteResult struct {
	gophercloud.ErrResult
}
