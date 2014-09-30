package tokens

import (
	"net/http"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
)

// RFC3339Milli describes the time format used by identity API responses.
const RFC3339Milli = "2006-01-02T15:04:05.999999Z"

// commonResult is the deferred result of a Create or a Get call.
type commonResult struct {
	gophercloud.CommonResult
	header http.Header
}

// Extract interprets a commonResult as a Token.
func (r commonResult) Extract() (*Token, error) {
	if r.Err != nil {
		return nil, r.Err
	}

	var response struct {
		Token struct {
			ExpiresAt string `mapstructure:"expires_at"`
		} `mapstructure:"token"`
	}

	var token Token

	// Parse the token itself from the stored headers.
	token.ID = r.header.Get("X-Subject-Token")

	err := mapstructure.Decode(r.Resp, &response)
	if err != nil {
		return nil, err
	}

	// Attempt to parse the timestamp.
	token.ExpiresAt, err = time.Parse(RFC3339Milli, response.Token.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// CreateResult is the deferred response from a Create call.
type CreateResult struct {
	commonResult
}

// createErr quickly creates a CreateResult that reports an error.
func createErr(err error) CreateResult {
	return CreateResult{
		commonResult: commonResult{
			CommonResult: gophercloud.CommonResult{Err: err},
			header:       nil,
		},
	}
}

// GetResult is the deferred response from a Get call.
type GetResult struct {
	commonResult
}

// Token is a string that grants a user access to a controlled set of services in an OpenStack provider.
// Each Token is valid for a set length of time.
type Token struct {
	// ID is the issued token.
	ID string

	// ExpiresAt is the timestamp at which this token will no longer be accepted.
	ExpiresAt time.Time
}
