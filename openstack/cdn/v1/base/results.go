package base

import (
	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// HomeDocument is a resource that contains all the resources for the CDN API.
type HomeDocument map[string]interface{}

// GetResult represents the result of a Get operation.
type GetResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a home document resource.
func (r GetResult) Extract() (*HomeDocument, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var res struct {
		HomeDocument *HomeDocument `mapstructure:"resources"`
	}

	err := mapstructure.Decode(r.Body, &res)

	return res.HomeDocument, err
}

// PingResult represents the result of a Ping operation.
type PingResult struct {
	gophercloud.ErrResult
}
