package gophercloud

import (
	"github.com/racker/perigee"
)

type AuthOptions struct {
	Username, Password, TenantId string
}

type AuthContainer struct {
	Auth Auth `json:"auth"`
}

type Auth struct {
	PasswordCredentials PasswordCredentials `json:"passwordCredentials"`
	TenantId            string              `json:"tenantId,omitempty"`
}

type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Access encapsulates the API token and its relevant fields, as well as the
// services catalog that Identity API returns once authenticated.  You'll probably
// rarely use this record directly, unless you intend on marshalling or unmarshalling
// Identity API JSON records yourself.
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
		Results: &struct{
			Access **Access `json:"access"`
		}{
			&access,
		},
	})
	return access, err
}
