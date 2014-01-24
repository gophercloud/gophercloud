package identity

import (
	"github.com/racker/perigee"
)

type AuthResults map[string]interface{}

// AuthOptions lets anyone calling Authenticate() supply the required access credentials.
// At present, only Identity V2 API support exists; therefore, only Username, Password,
// and optionally, TenantId are provided.  If future Identity API versions become available,
// alternative fields unique to those versions may appear here.
type AuthOptions struct {
	// Endpoint specifies the HTTP endpoint offering the Identity V2 API.
	// Required.
	Endpoint string

	// Username is required if using Identity V2 API.
	// Consult with your provider's control panel to discover your
	// account's username.
	Username string

	// At most one of Password or ApiKey is required if using Identity V2 API.
	// Consult with your provider's control panel to discover your
	// account's preferred method of authentication.
	Password, ApiKey string

	// The TenantId field is optional for the Identity V2 API.
	TenantId string

	// The TenantName can be specified instead of the TenantId
	TenantName string

	// AllowReauth should be set to true if you grant permission for Gophercloud to cache
	// your credentials in memory, and to allow Gophercloud to attempt to re-authenticate
	// automatically if/when your token expires.  If you set it to false, it will not cache
	// these settings, but re-authentication will not be possible.  This setting defaults
	// to false.
	AllowReauth bool
}

func Authenticate(options AuthOptions) (AuthResults, error) {
	var ar AuthResults

	if options.Endpoint == "" {
		return nil, ErrEndpoint
	}

	if (options.Username == "") || (options.Password == "" && options.ApiKey == "") {
		return nil, ErrCredentials
	}

	err := perigee.Post(options.Endpoint, perigee.Options{
		ReqBody: &AuthContainer{
			Auth: getAuthCredentials(options),
		},
		Results: &ar,
	})
	return ar, err
}

func getAuthCredentials(options AuthOptions) Auth {
	if options.ApiKey == "" {
		return Auth{
			PasswordCredentials: &struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}{
				Username: options.Username,
				Password: options.Password,
			},
			TenantId:   options.TenantId,
			TenantName: options.TenantName,
		}
	} else {
		return Auth{
			ApiKeyCredentials: &struct {
				Username string `json:"username"`
				ApiKey   string `json:"apiKey"`
			}{
				Username: options.Username,
				ApiKey:   options.ApiKey,
			},
			TenantId:   options.TenantId,
			TenantName: options.TenantName,
		}
	}
}

// AuthContainer provides a JSON encoding wrapper for passing credentials to the Identity
// service.  You will not work with this structure directly.
type AuthContainer struct {
	Auth Auth `json:"auth"`
}

// Auth provides a JSON encoding wrapper for passing credentials to the Identity
// service.  You will not work with this structure directly.
type Auth struct {
	PasswordCredentials interface{} `json:"passwordCredentials,omitempty"`
	ApiKeyCredentials   interface{} `json:"RAX-KSKEY:apiKeyCredentials,omitempty"`
	TenantId            string      `json:"tenantId,omitempty"`
	TenantName          string      `json:"tenantName,omitempty"`
}
