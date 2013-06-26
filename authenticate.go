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
	TenantId string `json:"tenantId,omitempty"`
}

type PasswordCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ProviderAccess interface {
	// ...
}

func (c *Context) Authenticate(provider string, options AuthOptions) (ProviderAccess, error) {
	p, err := c.ProviderByName(provider)
	if err != nil {
		return nil, err
	}
	if (options.Username == "") || (options.Password == "") {
		return nil, ErrCredentials
	}

	err = perigee.Post(p.AuthEndpoint, perigee.Options{
		CustomClient: c.httpClient,
		ReqBody:      &AuthContainer{
			Auth: Auth{
				PasswordCredentials: PasswordCredentials{
					Username: options.Username,
					Password: options.Password,
				},
				TenantId: options.TenantId,
			},
		},
	})
	return nil, err

	// if err != nil {
	// 	return err
	// }

	// c.isAuthenticated = true
	// c.token = id.access.Access.Token.Id
	// c.expires = id.access.Access.Token.Expires
	// c.tenantId = id.access.Access.Token.Tenant.Id
	// c.tenantName = id.access.Access.Token.Tenant.Name
}
