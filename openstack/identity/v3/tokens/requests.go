package tokens

import "github.com/gophercloud/gophercloud"

// Scope allows a created token to be limited to a specific domain or project.
type Scope struct {
	ProjectID   string `json:"scope.project.id,omitempty" not:"ProjectName,DomainID,DomainName"`
	ProjectName string `json:"scope.project.name,omitempty"`
	DomainID    string `json:"scope.project.id,omitempty" not:"ProjectName,ProjectID,DomainName"`
	DomainName  string `json:"scope.project.id,omitempty"`
}

// AuthOptionsBuilder describes any argument that may be passed to the Create call.
type AuthOptionsBuilder interface {
	// ToTokenV3CreateMap assembles the Create request body, returning an error if parameters are
	// missing or inconsistent.
	ToTokenV3CreateMap(*Scope) (map[string]interface{}, error)
}

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

func (opts AuthOptions) ToTokenV3CreateMap(scope *Scope) (map[string]interface{}, error) {
	type domainReq struct {
		ID   *string `json:"id,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	type projectReq struct {
		Domain *domainReq `json:"domain,omitempty"`
		Name   *string    `json:"name,omitempty"`
		ID     *string    `json:"id,omitempty"`
	}

	type userReq struct {
		ID       *string    `json:"id,omitempty"`
		Name     *string    `json:"name,omitempty"`
		Password string     `json:"password"`
		Domain   *domainReq `json:"domain,omitempty"`
	}

	type passwordReq struct {
		User userReq `json:"user"`
	}

	type tokenReq struct {
		ID string `json:"id"`
	}

	type identityReq struct {
		Methods  []string     `json:"methods"`
		Password *passwordReq `json:"password,omitempty"`
		Token    *tokenReq    `json:"token,omitempty"`
	}

	type scopeReq struct {
		Domain  *domainReq  `json:"domain,omitempty"`
		Project *projectReq `json:"project,omitempty"`
	}

	type authReq struct {
		Identity identityReq `json:"identity"`
		Scope    *scopeReq   `json:"scope,omitempty"`
	}

	type request struct {
		Auth authReq `json:"auth"`
	}

	// Populate the request structure based on the provided arguments. Create and return an error
	// if insufficient or incompatible information is present.
	var req request

	// Test first for unrecognized arguments.
	if opts.TenantID != "" {
		return nil, ErrTenantIDProvided{}
	}
	if opts.TenantName != "" {
		return nil, ErrTenantNameProvided{}
	}

	if opts.Password == "" {
		if opts.TokenID != "" {
			// Because we aren't using password authentication, it's an error to also provide any of the user-based authentication
			// parameters.
			if opts.Username != "" {
				return nil, ErrUsernameWithToken{}
			}
			if opts.UserID != "" {
				return nil, ErrUserIDWithToken{}
			}
			if opts.DomainID != "" {
				return nil, ErrDomainIDWithToken{}
			}
			if opts.DomainName != "" {
				return nil, ErrDomainNameWithToken{}
			}

			// Configure the request for Token authentication.
			req.Auth.Identity.Methods = []string{"token"}
			req.Auth.Identity.Token = &tokenReq{
				ID: opts.TokenID,
			}
		} else {
			// If no password or token ID are available, authentication can't continue.
			return nil, ErrMissingPassword{}
		}
	} else {
		// Password authentication.
		req.Auth.Identity.Methods = []string{"password"}

		// At least one of Username and UserID must be specified.
		if opts.Username == "" && opts.UserID == "" {
			return nil, ErrUsernameOrUserID{}
		}

		if opts.Username != "" {
			// If Username is provided, UserID may not be provided.
			if opts.UserID != "" {
				return nil, ErrUsernameOrUserID{}
			}

			// Either DomainID or DomainName must also be specified.
			if opts.DomainID == "" && opts.DomainName == "" {
				return nil, ErrDomainIDOrDomainName{}
			}

			if opts.DomainID != "" {
				if opts.DomainName != "" {
					return nil, ErrDomainIDOrDomainName{}
				}

				// Configure the request for Username and Password authentication with a DomainID.
				req.Auth.Identity.Password = &passwordReq{
					User: userReq{
						Name:     &opts.Username,
						Password: opts.Password,
						Domain:   &domainReq{ID: &opts.DomainID},
					},
				}
			}

			if opts.DomainName != "" {
				// Configure the request for Username and Password authentication with a DomainName.
				req.Auth.Identity.Password = &passwordReq{
					User: userReq{
						Name:     &opts.Username,
						Password: opts.Password,
						Domain:   &domainReq{Name: &opts.DomainName},
					},
				}
			}
		}

		if opts.UserID != "" {
			// If UserID is specified, neither DomainID nor DomainName may be.
			if opts.DomainID != "" {
				return nil, ErrDomainIDWithUserID{}
			}
			if opts.DomainName != "" {
				return nil, ErrDomainNameWithUserID{}
			}

			// Configure the request for UserID and Password authentication.
			req.Auth.Identity.Password = &passwordReq{
				User: userReq{ID: &opts.UserID, Password: opts.Password},
			}
		}
	}

	// Add a "scope" element if a Scope has been provided.
	if scope != nil {
		if scope.ProjectName != "" {
			// ProjectName provided: either DomainID or DomainName must also be supplied.
			// ProjectID may not be supplied.
			if scope.DomainID == "" && scope.DomainName == "" {
				return nil, ErrScopeDomainIDOrDomainName{}
			}
			if scope.ProjectID != "" {
				return nil, ErrScopeProjectIDOrProjectName{}
			}

			if scope.DomainID != "" {
				// ProjectName + DomainID
				req.Auth.Scope = &scopeReq{
					Project: &projectReq{
						Name:   &scope.ProjectName,
						Domain: &domainReq{ID: &scope.DomainID},
					},
				}
			}

			if scope.DomainName != "" {
				// ProjectName + DomainName
				req.Auth.Scope = &scopeReq{
					Project: &projectReq{
						Name:   &scope.ProjectName,
						Domain: &domainReq{Name: &scope.DomainName},
					},
				}
			}
		} else if scope.ProjectID != "" {
			// ProjectID provided. ProjectName, DomainID, and DomainName may not be provided.
			if scope.DomainID != "" {
				return nil, ErrScopeProjectIDAlone{}
			}
			if scope.DomainName != "" {
				return nil, ErrScopeProjectIDAlone{}
			}

			// ProjectID
			req.Auth.Scope = &scopeReq{
				Project: &projectReq{ID: &scope.ProjectID},
			}
		} else if scope.DomainID != "" {
			// DomainID provided. ProjectID, ProjectName, and DomainName may not be provided.
			if scope.DomainName != "" {
				return nil, ErrScopeDomainIDOrDomainName{}
			}

			// DomainID
			req.Auth.Scope = &scopeReq{
				Domain: &domainReq{ID: &scope.DomainID},
			}
		} else if scope.DomainName != "" {
			return nil, ErrScopeDomainName{}
		} else {
			return nil, ErrScopeEmpty{}
		}
	}

	b, err2 := gophercloud.BuildRequestBody(req, "")
	if err2 != nil {
		return nil, err2
	}
	return b, nil
}

func subjectTokenHeaders(c *gophercloud.ServiceClient, subjectToken string) map[string]string {
	return map[string]string{
		"X-Subject-Token": subjectToken,
	}
}

// Create authenticates and either generates a new token, or changes the Scope of an existing token.
func Create(c *gophercloud.ServiceClient, opts AuthOptionsBuilder, scopeOpts *Scope) (r CreateResult) {
	b, err := opts.ToTokenV3CreateMap(scopeOpts)
	if err != nil {
		r.Err = err
		return
	}
	resp, err := c.Post(tokenURL(c), b, &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{"X-Auth-Token": ""},
	})
	if resp != nil {
		r.Err = err
		r.Header = resp.Header
	}
	return
}

// Get validates and retrieves information about another token.
func Get(c *gophercloud.ServiceClient, token string) (r GetResult) {
	resp, err := c.Get(tokenURL(c), &r.Body, &gophercloud.RequestOpts{
		MoreHeaders: subjectTokenHeaders(c, token),
		OkCodes:     []int{200, 203},
	})
	if resp != nil {
		r.Err = err
		r.Header = resp.Header
	}
	return
}

// Validate determines if a specified token is valid or not.
func Validate(c *gophercloud.ServiceClient, token string) (bool, error) {
	resp, err := c.Request("HEAD", tokenURL(c), &gophercloud.RequestOpts{
		MoreHeaders: subjectTokenHeaders(c, token),
		OkCodes:     []int{204, 404},
	})
	if err != nil {
		return false, err
	}

	return resp.StatusCode == 204, nil
}

// Revoke immediately makes specified token invalid.
func Revoke(c *gophercloud.ServiceClient, token string) (r RevokeResult) {
	_, r.Err = c.Delete(tokenURL(c), &gophercloud.RequestOpts{
		MoreHeaders: subjectTokenHeaders(c, token),
	})
	return
}
