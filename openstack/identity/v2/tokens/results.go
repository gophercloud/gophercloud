package tokens

import (
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tenants"
)

// Token provides only the most basic information related to an authentication token.
type Token struct {
	// ID provides the primary means of identifying a user to the OpenStack API.
	// OpenStack defines this field as an opaque value, so do not depend on its content.
	// It is safe, however, to compare for equality.
	ID string

	// ExpiresAt provides a timestamp in ISO 8601 format, indicating when the authentication token becomes invalid.
	// After this point in time, future API requests made using this authentication token will respond with errors.
	// Either the caller will need to reauthenticate manually, or more preferably, the caller should exploit automatic re-authentication.
	// See the AuthOptions structure for more details.
	ExpiresAt time.Time

	// Tenant provides information about the tenant to which this token grants access.
	Tenant tenants.Tenant
}

// CreateResult defers the interpretation of a created token.
// Use ExtractToken() to interpret it as a Token, or ExtractServiceCatalog() to interpret it as a service catalog.
type CreateResult struct {
	gophercloud.CommonResult
}

// ExtractToken returns the just-created Token from a CreateResult.
func (result CreateResult) ExtractToken() (*Token, error) {
	var response struct {
		Access struct {
			Token struct {
				Expires string         `mapstructure:"expires"`
				ID      string         `mapstructure:"id"`
				Tenant  tenants.Tenant `mapstructure:"tenant"`
			} `mapstructure:"token"`
		} `mapstructure:"access"`
	}

	err := mapstructure.Decode(result.Resp, &response)
	if err != nil {
		return nil, err
	}

	expiresTs, err := time.Parse(gophercloud.RFC3339Milli, response.Access.Token.Expires)
	if err != nil {
		return nil, err
	}

	return &Token{
		ID:        response.Access.Token.ID,
		ExpiresAt: expiresTs,
		Tenant:    response.Access.Token.Tenant,
	}, nil
}

// createErr quickly packs an error in a CreateResult.
func createErr(err error) CreateResult {
	return CreateResult{gophercloud.CommonResult{Err: err}}
}
