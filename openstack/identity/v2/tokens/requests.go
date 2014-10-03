package tokens

import (
	"github.com/racker/perigee"
	"github.com/rackspace/gophercloud"
)

// Create authenticates to the identity service and attempts to acquire a Token.
// If successful, the CreateResult
// Generally, rather than interact with this call directly, end users should call openstack.AuthenticatedClient(),
// which abstracts all of the gory details about navigating service catalogs and such.
func Create(client *gophercloud.ServiceClient, auth gophercloud.AuthOptions) CreateResult {
	type passwordCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type apiKeyCredentials struct {
		Username string `json:"username"`
		APIKey   string `json:"apiKey"`
	}

	var request struct {
		Auth struct {
			PasswordCredentials *passwordCredentials `json:"passwordCredentials,omitempty"`
			APIKeyCredentials   *apiKeyCredentials   `json:"RAX-KSKEY:apiKeyCredentials,omitempty"`
			TenantID            string               `json:"tenantId,omitempty"`
			TenantName          string               `json:"tenantName,omitempty"`
		} `json:"auth"`
	}

	// Error out if an unsupported auth option is present.
	if auth.UserID != "" {
		return createErr(ErrUserIDProvided)
	}
	if auth.DomainID != "" {
		return createErr(ErrDomainIDProvided)
	}
	if auth.DomainName != "" {
		return createErr(ErrDomainNameProvided)
	}

	// Username is always required.
	if auth.Username == "" {
		return createErr(ErrUsernameRequired)
	}

	// Populate either PasswordCredentials or APIKeyCredentials
	if auth.Password != "" {
		if auth.APIKey != "" {
			return createErr(ErrPasswordOrAPIKey)
		}

		// Username + Password
		request.Auth.PasswordCredentials = &passwordCredentials{
			Username: auth.Username,
			Password: auth.Password,
		}
	} else if auth.APIKey != "" {
		// API key authentication.
		request.Auth.APIKeyCredentials = &apiKeyCredentials{
			Username: auth.Username,
			APIKey:   auth.APIKey,
		}
	} else {
		return createErr(ErrPasswordOrAPIKey)
	}

	// Populate the TenantName or TenantID, if provided.
	request.Auth.TenantID = auth.TenantID
	request.Auth.TenantName = auth.TenantName

	var result CreateResult
	_, result.Err = perigee.Request("POST", listURL(client), perigee.Options{
		ReqBody: &request,
		Results: &result.Resp,
		OkCodes: []int{200, 203},
	})
	return result
}
