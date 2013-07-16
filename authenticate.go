package gophercloud

import (
	"github.com/racker/perigee"
)

// AuthOptions lets anyone calling Authenticate() supply the required access credentials.
// At present, only Identity V2 API support exists; therefore, only Username, Password,
// and optionally, TenantId are provided.  If future Identity API versions become available,
// alternative fields unique to those versions may appear here.
type AuthOptions struct {
	// Username and Password are required if using Identity V2 API.
	// Consult with your provider's control panel to discover your
	// account's username and password.
	Username, Password string

	// The TenantId field is optional for the Identity V2 API.
	TenantId string
}

// AuthContainer provides a JSON encoding wrapper for passing credentials to the Identity
// service.  You will not work with this structure directly.
type AuthContainer struct {
	Auth Auth `json:"auth"`
}

// Auth provides a JSON encoding wrapper for passing credentials to the Identity
// service.  You will not work with this structure directly.
type Auth struct {
	PasswordCredentials PasswordCredentials `json:"passwordCredentials"`
	TenantId            string              `json:"tenantId,omitempty"`
}

// PasswordCredentials provides a JSON encoding wrapper for passing credentials to the Identity
// service.  You will not work with this structure directly.
type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Access encapsulates the API token and its relevant fields, as well as the
// services catalog that Identity API returns once authenticated.
type Access struct {
	Token          Token
	ServiceCatalog []CatalogEntry
	User           User
}

// Token encapsulates an authentication token and when it expires.  It also includes
// tenant information if available.
type Token struct {
	Id, Expires string
	Tenant      Tenant
}

// Tenant encapsulates tenant authentication information.  If, after authentication,
// no tenant information is supplied, both Id and Name will be "".
type Tenant struct {
	Id, Name string
}

// User encapsulates the user credentials, and provides visibility in what
// the user can do through its role assignments.
type User struct {
	Id, Name          string
	XRaxDefaultRegion string `json:"RAX-AUTH:defaultRegion"`
	Roles             []Role
}

// Role encapsulates a permission that a user can rely on.
type Role struct {
	Description, Id, Name string
}

// CatalogEntry encapsulates a service catalog record.
type CatalogEntry struct {
	Name, Type string
	Endpoints  []EntryEndpoint
}

// EntryEndpoint encapsulates how to get to the API of some service.
type EntryEndpoint struct {
	Region, TenantId                    string
	PublicURL, InternalURL              string
	VersionId, VersionInfo, VersionList string
}

// Authenticate() grants access to the OpenStack-compatible provider API.
//
// Providers are identified through a unique key string.
// See the RegisterProvider() method for more details.
//
// The supplied AuthOptions instance allows the client to specify only those credentials
// relevant for the authentication request.  At present, support exists for OpenStack
// Identity V2 API only; support for V3 will become available as soon as documentation for it
// becomes readily available.
//
// For Identity V2 API requirements, you must provide at least the Username and Password
// options.  The TenantId field is optional, and defaults to "".
func (c *Context) Authenticate(provider string, options AuthOptions) (*Access, error) {
	var access *Access

	p, err := c.ProviderByName(provider)
	if err != nil {
		return nil, err
	}
	if (options.Username == "") || (options.Password == "") {
		return nil, ErrCredentials
	}

	err = perigee.Post(p.AuthEndpoint, perigee.Options{
		CustomClient: c.httpClient,
		ReqBody: &AuthContainer{
			Auth: Auth{
				PasswordCredentials: PasswordCredentials{
					Username: options.Username,
					Password: options.Password,
				},
				TenantId: options.TenantId,
			},
		},
		Results: &struct {
			Access **Access `json:"access"`
		}{
			&access,
		},
	})
	return access, err
}

// See AccessProvider interface definition for details.
func (a *Access) FirstEndpointUrlByCriteria(ac ApiCriteria) string {
	ep := FindFirstEndpointByCriteria(a.ServiceCatalog, ac)
	urls := []string{ep.PublicURL, ep.InternalURL}
	return urls[ac.UrlChoice]
}

// See AccessProvider interface definition for details.
func (a *Access) AuthToken() string {
	return a.Token.Id
}

// See AccessProvider interface definition for details.
func (a *Access) Revoke(tok string) error {
	return nil
}
