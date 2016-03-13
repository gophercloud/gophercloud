package gophercloud

/*
type AuthOptionsBuilder interface {
	ToTokenCreateMap() (map[string]interface{}, error)
}
*/

/*
AuthOptions stores information needed to authenticate to an OpenStack cluster.
You can populate one manually, or use a provider's AuthOptionsFromEnv() function
to read relevant information from the standard environment variables. Pass one
to a provider's AuthenticatedClient function to authenticate and obtain a
ProviderClient representing an active session on that provider.

Its fields are the union of those recognized by each identity implementation and
provider.
*/
type AuthOptions struct {
	// IdentityEndpoint specifies the HTTP endpoint that is required to work with
	// the Identity API of the appropriate version. While it's ultimately needed by
	// all of the identity services, it will often be populated by a provider-level
	// function.
	IdentityEndpoint string `json:"-"`

	// Username is required if using Identity V2 API. Consult with your provider's
	// control panel to discover your account's username. In Identity V3, either
	// UserID or a combination of Username and DomainID or DomainName are needed.
	Username string `json:"username,omitempty"`
	UserID   string `json:"id,omitempty"`

	Password string `json:"password,omitempty"`

	// At most one of DomainID and DomainName must be provided if using Username
	// with Identity V3. Otherwise, either are optional.
	DomainID   string `json:"id,omitempty"`
	DomainName string `json:"name,omitempty"`

	// The TenantID and TenantName fields are optional for the Identity V2 API.
	// Some providers allow you to specify a TenantName instead of the TenantId.
	// Some require both. Your provider's authentication policies will determine
	// how these fields influence authentication.
	TenantID   string `json:"tenantId,omitempty"`
	TenantName string `json:"tenantName,omitempty"`

	// AllowReauth should be set to true if you grant permission for Gophercloud to
	// cache your credentials in memory, and to allow Gophercloud to attempt to
	// re-authenticate automatically if/when your token expires.  If you set it to
	// false, it will not cache these settings, but re-authentication will not be
	// possible.  This setting defaults to false.
	AllowReauth bool `json:"-"`

	// TokenID allows users to authenticate (possibly as another user) with an
	// authentication token ID.
	TokenID string
}

// ToTokenV2CreateMap allows AuthOptions to satisfy the AuthOptionsBuilder
// interface in the v2 tokens package
func (opts AuthOptions) ToTokenV2CreateMap() (map[string]interface{}, error) {
	v2Opts := AuthOptionsV2{
		PasswordCredentials: &PasswordCredentialsV2{
			Username: opts.Username,
			Password: opts.Password,
		},
		TenantID:   opts.TenantID,
		TenantName: opts.TenantName,
		TokenCredentials: &TokenCredentialsV2{
			ID: opts.TokenID,
		},
	}

	b, err := BuildRequestBody(v2Opts, "auth")
	if err != nil {
		return nil, err
	}
	/*
		if opts.TokenID == "" {
			delete(b["auth"].(map[string]interface{}), "token")
			return b, nil
		}

		delete(b["auth"].(map[string]interface{}), "passwordCredentials")*/
	return b, nil
}

func (opts AuthOptions) ToTokenV3CreateMap(scope *ScopeOptsV3) (map[string]interface{}, error) {
	var methods []string
	if opts.TokenID != "" {
		methods = []string{"token"}
	} else {
		methods = []string{"password"}
	}

	v3Opts := AuthOptionsV3{
		Identity: &IdentityCredentialsV3{
			Methods: methods,
			PasswordCredentials: &PasswordCredentialsV3{
				User: &UserV3{
					ID:       opts.UserID,
					Name:     opts.Username,
					Password: opts.Password,
					Domain: &DomainV3{
						ID:   opts.DomainID,
						Name: opts.DomainName,
					},
				},
			},
			TokenCredentials: &TokenCredentialsV3{
				ID: opts.TokenID,
			},
		},
	}

	if scope != nil {
		v3Opts.Scope = &ScopeV3{
			Domain: &ScopeDomainV3{
				ID:   scope.DomainID,
				Name: scope.DomainName,
			},
			Project: &ScopeProjectV3{
				Domain: &ScopeProjectDomainV3{
					ID:   scope.DomainID,
					Name: scope.DomainName,
				},
				ID:   scope.ProjectID,
				Name: scope.ProjectName,
			},
		}
	}

	b, err := BuildRequestBody(v3Opts, "auth")
	if err != nil {
		return nil, err
	}
	/*
		if opts.TokenID == "" {
			delete(b["auth"].(map[string]interface{}), "token")
			return b, nil
		}

		delete(b["auth"].(map[string]interface{}), "passwordCredentials")*/
	return b, nil
}
