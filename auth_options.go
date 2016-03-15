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
	/*
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
		if options.APIKey != "" {
			return createErr(ErrAPIKeyProvided)
		}
		if options.TenantID != "" {
			return createErr(ErrTenantIDProvided)
		}
		if options.TenantName != "" {
			return createErr(ErrTenantNameProvided)
		}

		if options.Password == "" {
			if c.TokenID != "" {
				// Because we aren't using password authentication, it's an error to also provide any of the user-based authentication
				// parameters.
				if options.Username != "" {
					return createErr(ErrUsernameWithToken)
				}
				if options.UserID != "" {
					return createErr(ErrUserIDWithToken)
				}
				if options.DomainID != "" {
					return createErr(ErrDomainIDWithToken)
				}
				if options.DomainName != "" {
					return createErr(ErrDomainNameWithToken)
				}

				// Configure the request for Token authentication.
				req.Auth.Identity.Methods = []string{"token"}
				req.Auth.Identity.Token = &tokenReq{
					ID: c.TokenID,
				}
			} else {
				// If no password or token ID are available, authentication can't continue.
				return createErr(ErrMissingPassword)
			}
		} else {
			// Password authentication.
			req.Auth.Identity.Methods = []string{"password"}

			// At least one of Username and UserID must be specified.
			if options.Username == "" && options.UserID == "" {
				return createErr(ErrUsernameOrUserID)
			}

			if options.Username != "" {
				// If Username is provided, UserID may not be provided.
				if options.UserID != "" {
					return createErr(ErrUsernameOrUserID)
				}

				// Either DomainID or DomainName must also be specified.
				if options.DomainID == "" && options.DomainName == "" {
					return createErr(ErrDomainIDOrDomainName)
				}

				if options.DomainID != "" {
					if options.DomainName != "" {
						return createErr(ErrDomainIDOrDomainName)
					}

					// Configure the request for Username and Password authentication with a DomainID.
					req.Auth.Identity.Password = &passwordReq{
						User: userReq{
							Name:     &options.Username,
							Password: options.Password,
							Domain:   &domainReq{ID: &options.DomainID},
						},
					}
				}

				if options.DomainName != "" {
					// Configure the request for Username and Password authentication with a DomainName.
					req.Auth.Identity.Password = &passwordReq{
						User: userReq{
							Name:     &options.Username,
							Password: options.Password,
							Domain:   &domainReq{Name: &options.DomainName},
						},
					}
				}
			}

			if options.UserID != "" {
				// If UserID is specified, neither DomainID nor DomainName may be.
				if options.DomainID != "" {
					return createErr(ErrDomainIDWithUserID)
				}
				if options.DomainName != "" {
					return createErr(ErrDomainNameWithUserID)
				}

				// Configure the request for UserID and Password authentication.
				req.Auth.Identity.Password = &passwordReq{
					User: userReq{ID: &options.UserID, Password: options.Password},
				}
			}
		}

		// Add a "scope" element if a Scope has been provided.
		if scope != nil {
			if scope.ProjectName != "" {
				// ProjectName provided: either DomainID or DomainName must also be supplied.
				// ProjectID may not be supplied.
				if scope.DomainID == "" && scope.DomainName == "" {
					return createErr(ErrScopeDomainIDOrDomainName)
				}
				if scope.ProjectID != "" {
					return createErr(ErrScopeProjectIDOrProjectName)
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
					return createErr(ErrScopeProjectIDAlone)
				}
				if scope.DomainName != "" {
					return createErr(ErrScopeProjectIDAlone)
				}

				// ProjectID
				req.Auth.Scope = &scopeReq{
					Project: &projectReq{ID: &scope.ProjectID},
				}
			} else if scope.DomainID != "" {
				// DomainID provided. ProjectID, ProjectName, and DomainName may not be provided.
				if scope.DomainName != "" {
					return createErr(ErrScopeDomainIDOrDomainName)
				}

				// DomainID
				req.Auth.Scope = &scopeReq{
					Domain: &domainReq{ID: &scope.DomainID},
				}
			} else if scope.DomainName != "" {
				return createErr(ErrScopeDomainName)
			} else {
				return createErr(ErrScopeEmpty)
			}
		}
	*/

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

	return b, nil
}
