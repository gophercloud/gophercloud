package identity

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// AuthResults encapsulates the raw results from an authentication request.
// As OpenStack allows extensions to influence the structure returned in
// ways that Gophercloud cannot predict at compile-time, you should use
// type-safe accessors to work with the data represented by this type,
// such as ServiceCatalog() and Token().
type AuthResults map[string]interface{}

// Authenticate passes the supplied credentials to the OpenStack provider for authentication.
// If successful, the caller may use Token() to retrieve the authentication token,
// and ServiceCatalog() to retrieve the set of services available to the API user.
func Authenticate(service gophercloud.ServiceClient, options gophercloud.AuthOptions) (AuthResults, error) {
	type AuthContainer struct {
		Auth auth `json:"auth"`
	}

	var ar AuthResults

	if options.IdentityEndpoint == "" {
		return nil, ErrEndpoint
	}

	if (options.Username == "") || (options.Password == "" && options.APIKey == "") {
		return nil, ErrCredentials
	}

	url := options.IdentityEndpoint + "/tokens"
	err := perigee.Post(url, perigee.Options{
		ReqBody: &AuthContainer{
			Auth: getAuthCredentials(options),
		},
		Results: &ar,
	})
	return ar, err
}

func getAuthCredentials(options gophercloud.AuthOptions) auth {
	if options.APIKey == "" {
		return auth{
			PasswordCredentials: &struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}{
				Username: options.Username,
				Password: options.Password,
			},
			TenantID:   options.TenantID,
			TenantName: options.TenantName,
		}
	}
	return auth{
		APIKeyCredentials: &struct {
			Username string `json:"username"`
			APIKey   string `json:"apiKey"`
		}{
			Username: options.Username,
			APIKey:   options.APIKey,
		},
		TenantID:   options.TenantID,
		TenantName: options.TenantName,
	}
}

type auth struct {
	PasswordCredentials interface{} `json:"passwordCredentials,omitempty"`
	APIKeyCredentials   interface{} `json:"RAX-KSKEY:apiKeyCredentials,omitempty"`
	TenantID            string      `json:"tenantId,omitempty"`
	TenantName          string      `json:"tenantName,omitempty"`
}

// GetExtensions returns the OpenStack extensions available from this service.
func GetExtensions(options gophercloud.AuthOptions) (ExtensionsResult, error) {
	var exts ExtensionsResult

	url := options.IdentityEndpoint + "/extensions"
	err := perigee.Get(url, perigee.Options{
		Results: &exts,
	})
	return exts, err
}
