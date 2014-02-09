package identity

import (
	"github.com/racker/perigee"
)

// AuthResults encapsulates the raw results from an authentication request.
// As OpenStack allows extensions to influence the structure returned in
// ways that Gophercloud cannot predict at compile-time, you should use
// type-safe accessors to work with the data represented by this type,
// such as ServiceCatalog() and Token().
type AuthResults map[string]interface{}

// AuthOptions lets anyone calling Authenticate() supply the required access
// credentials.  At present, only Identity V2 API support exists; therefore,
// only Username, Password, and optionally, TenantId are provided.  If future
// Identity API versions become available, alternative fields unique to those
// versions may appear here.
//
// Endpoint specifies the HTTP endpoint offering the Identity V2 API.
// Required.
//
// Username is required if using Identity V2 API.  Consult with your provider's
// control panel to discover your account's username.
//
// At most one of Password or ApiKey is required if using Identity V2 API.
// Consult with your provider's control panel to discover your account's
// preferred method of authentication.
//
// The TenantId and TenantName fields are optional for the Identity V2 API.
// Some providers allow you to specify a TenantName instead of the TenantId.
// Some require both.  Your provider's authentication policies will determine
// how these fields influence authentication.
//
// AllowReauth should be set to true if you grant permission for Gophercloud to
// cache your credentials in memory, and to allow Gophercloud to attempt to
// re-authenticate automatically if/when your token expires.  If you set it to
// false, it will not cache these settings, but re-authentication will not be
// possible.  This setting defaults to false.
type AuthOptions struct {
	Endpoint         string
	Username         string
	Password, ApiKey string
	TenantId         string
	TenantName       string
	AllowReauth      bool
}

// Authenticate passes the supplied credentials to the OpenStack provider for authentication.
// If successful, the caller may use Token() to retrieve the authentication token,
// and ServiceCatalog() to retrieve the set of services available to the API user.
func Authenticate(options AuthOptions) (AuthResults, error) {
	type AuthContainer struct {
		Auth auth `json:"auth"`
	}

	var ar AuthResults

	if options.Endpoint == "" {
		return nil, ErrEndpoint
	}

	if (options.Username == "") || (options.Password == "" && options.ApiKey == "") {
		return nil, ErrCredentials
	}

	url := options.Endpoint + "/tokens"
	err := perigee.Post(url, perigee.Options{
		ReqBody: &AuthContainer{
			Auth: getAuthCredentials(options),
		},
		Results: &ar,
	})
	return ar, err
}

func getAuthCredentials(options AuthOptions) auth {
	if options.ApiKey == "" {
		return auth{
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
		return auth{
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

type auth struct {
	PasswordCredentials interface{} `json:"passwordCredentials,omitempty"`
	ApiKeyCredentials   interface{} `json:"RAX-KSKEY:apiKeyCredentials,omitempty"`
	TenantId            string      `json:"tenantId,omitempty"`
	TenantName          string      `json:"tenantName,omitempty"`
}

func GetExtensions(options AuthOptions) (ExtensionsResult, error) {
	var exts ExtensionsResult

	url := options.Endpoint + "/extensions"
	err := perigee.Get(url, perigee.Options{
		Results: &exts,
	})
	return exts, err
}
