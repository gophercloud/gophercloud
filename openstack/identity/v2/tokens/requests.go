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

	var request struct {
		Auth struct {
			PasswordCredentials *passwordCredentials `json:"passwordCredentials"`
			TenantID            string               `json:"tenantId,omitempty"`
			TenantName          string               `json:"tenantName,omitempty"`
		} `json:"auth"`
	}

	// Error out if an unsupported auth option is present.
	if auth.UserID != "" {
		return createErr(ErrUserIDProvided)
	}
	if auth.APIKey != "" {
		return createErr(ErrAPIKeyProvided)
	}
	if auth.DomainID != "" {
		return createErr(ErrDomainIDProvided)
	}
	if auth.DomainName != "" {
		return createErr(ErrDomainNameProvided)
	}

	// Username and Password are always required.
	if auth.Username == "" {
		return createErr(ErrUsernameRequired)
	}
	if auth.Password == "" {
		return createErr(ErrPasswordRequired)
	}

	// Populate the request.
	request.Auth.PasswordCredentials = &passwordCredentials{
		Username: auth.Username,
		Password: auth.Password,
	}
	request.Auth.TenantID = auth.TenantID
	request.Auth.TenantName = auth.TenantName

	var result CreateResult
	_, result.Err = perigee.Request("POST", CreateURL(client), perigee.Options{
		ReqBody: &request,
		Results: &result.Resp,
		OkCodes: []int{200, 203},
	})
	return result
}
