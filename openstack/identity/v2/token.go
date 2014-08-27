package identity

import (
	"github.com/mitchellh/mapstructure"
)

// Token provides only the most basic information related to an authentication token.
//
// Id provides the primary means of identifying a user to the OpenStack API.
// OpenStack defines this field as an opaque value, so do not depend on its content.
// It is safe, however, to compare for equality.
//
// Expires provides a timestamp in ISO 8601 format, indicating when the authentication token becomes invalid.
// After this point in time, future API requests made using this authentication token will respond with errors.
// Either the caller will need to reauthenticate manually, or more preferably, the caller should exploit automatic re-authentication.
// See the AuthOptions structure for more details.
//
// TenantId provides the canonical means of identifying a tenant.
// As with Id, this field is defined to be opaque, so do not depend on its content.
// It is safe, however, to compare for equality.
//
// TenantName provides a human-readable tenant name corresponding to the TenantId.
type Token struct {
	Id, Expires          string
	TenantId, TenantName string
}

// GetToken, if successful, yields an unpacked collection of fields related to the user's access credentials, called a "token."
// See the Token structure for more details.
func GetToken(m AuthResults) (*Token, error) {
	type (
		Tenant struct {
			Id   string
			Name string
		}

		TokenDesc struct {
			Id      string `mapstructure:"id"`
			Expires string `mapstructure:"expires"`
			Tenant
		}
	)

	accessMap, err := getSubmap(m, "access")
	if err != nil {
		return nil, err
	}
	tokenMap, err := getSubmap(accessMap, "token")
	if err != nil {
		return nil, err
	}
	t := &TokenDesc{}
	err = mapstructure.Decode(tokenMap, t)
	if err != nil {
		return nil, err
	}
	td := &Token{
		Id:         t.Id,
		Expires:    t.Expires,
		TenantId:   t.Tenant.Id,
		TenantName: t.Tenant.Name,
	}
	return td, nil
}

func getSubmap(m map[string]interface{}, name string) (map[string]interface{}, error) {
	entry, ok := m[name]
	if !ok {
		return nil, ErrNotImplemented
	}
	return entry.(map[string]interface{}), nil
}
